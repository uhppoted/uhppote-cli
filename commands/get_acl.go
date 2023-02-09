package commands

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/uhppoted/uhppoted-lib/acl"
	"github.com/uhppoted/uhppoted-lib/config"
	"os"
	"strings"
)

var GetACLCmd = GetACL{}

type GetACL struct {
}

func (c *GetACL) Execute(ctx Context) error {
	if ctx.config == nil {
		return fmt.Errorf("get-acl requires a valid configuration file")
	}

	list, errors := acl.GetACL(ctx.uhppote, ctx.devices)
	if len(errors) > 0 {
		return fmt.Errorf("%v", errors)
	}

	for k, l := range list {
		fmt.Printf("   ... %v  Retrieved %v records\n", k, len(l))
	}

	file, err := c.getACLFile()
	if err != nil {
		return err
	}

	if file != "" {
		var w bytes.Buffer
		if err = acl.MakeTSV(list, ctx.devices, &w); err != nil {
			return err
		}

		return os.WriteFile(file, w.Bytes(), 0660)
	}

	var w strings.Builder
	if err = acl.MakeFlatFile(list, ctx.devices, &w); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(w.String())
	fmt.Println()

	return nil
}

func (c *GetACL) getACLFile() (string, error) {
	if len(flag.Args()) < 2 {
		return "", nil
	}

	file := flag.Arg(1)
	stat, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return file, nil
		}
		return "", err
	}

	if stat.Mode().IsDir() {
		return "", fmt.Errorf("file '%s' is a directory", file)
	}

	if !stat.Mode().IsRegular() {
		return "", fmt.Errorf("file '%s' is not a real file", file)
	}

	return file, nil
}

func (c *GetACL) CLI() string {
	return "get-acl"
}

func (c *GetACL) Description() string {
	return "Retrieves the access control list as a TSV file for the configured controllers"
}

func (c *GetACL) Usage() string {
	return "<file>"
}

func (c *GetACL) Help() {
	fmt.Println("Usage: uhppote-cli [options] get-acl <TSV file>")
	fmt.Println()
	fmt.Println(" Retrieves the cards from the access controllers defined in the configuration file, reformats as")
	fmt.Println(" an access control list and writes to the specified TSV file")
	fmt.Println()
	fmt.Println("  <TSV file>  (optional) file to write TSV access control list. Writes to stdout if not provided")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --config  File path for the 'conf' file containing the controller configuration")
	fmt.Printf("              (defaults to %s)\n", config.DefaultConfig)
	fmt.Println("    --debug   Displays internal information for diagnosing errors")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli --debug --config .config get-acl \"uhppote.tsv\"")
	fmt.Println()
}

// Returns true - configuration is not optional for this command to return valid information.
func (c *GetACL) RequiresConfig() bool {
	return true
}
