package commands

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/uhppoted/uhppoted-api/acl"
	"io/ioutil"
	"os"
	"strings"
)

type GetACL struct {
}

func (c *GetACL) Execute(ctx Context) error {
	if ctx.config == nil {
		return errors.New("get-acl requires a valid configuration file")
	}

	err := ctx.config.Verify()
	if err != nil {
		return err
	}

	devices := getDevices(&ctx)
	list, err := acl.GetACL(ctx.uhppote, devices)
	if err != nil {
		return err
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
		if err = acl.MakeTSV(list, devices, &w); err != nil {
			return err
		}

		return ioutil.WriteFile(file, w.Bytes(), 0660)
	}

	var w strings.Builder
	if err = acl.MakeFlatFile(list, devices, &w); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(w.String())
	fmt.Println()

	fmt.Println(w.String())
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
		return "", errors.New(fmt.Sprintf("File '%s' is a directory", file))
	}

	if !stat.Mode().IsRegular() {
		return "", errors.New(fmt.Sprintf("File '%s' is not a real file", file))
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
	fmt.Printf("              (defaults to %s)\n", DEFAULT_CONFIG)
	fmt.Println("    --debug   Displays vaguely useful internal information")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli --debug --config .config get-acl \"uhppote.tsv\"")
	fmt.Println()
}
