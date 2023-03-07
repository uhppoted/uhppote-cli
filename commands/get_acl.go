package commands

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/uhppoted/uhppote-core/uhppote"

	"github.com/uhppoted/uhppoted-lib/acl"
	"github.com/uhppoted/uhppoted-lib/config"
)

var GetACLCmd = GetACL{
	file:    "",
	withPIN: false,
}

type GetACL struct {
	file    string
	withPIN bool
}

func (c *GetACL) Execute(ctx Context) error {
	if ctx.config == nil {
		return fmt.Errorf("get-acl requires a valid configuration file")
	}

	if err := c.parseArgs(); err != nil {
		return err
	}

	tsv := func(list acl.ACL, devices []uhppote.Device, w io.Writer) error {
		if c.withPIN {
			return acl.MakeTSVWithPIN(list, devices, w)
		} else {
			return acl.MakeTSV(list, devices, w)
		}
	}

	txt := func(list acl.ACL, devices []uhppote.Device, w io.Writer) error {
		if c.withPIN {
			return acl.MakeFlatFileWithPIN(list, ctx.devices, w)
		} else {
			return acl.MakeFlatFile(list, ctx.devices, w)
		}
	}

	list, errors := acl.GetACL(ctx.uhppote, ctx.devices)
	if len(errors) > 0 {
		return fmt.Errorf("%v", errors)
	}

	for k, l := range list {
		fmt.Printf("   ... %v  Retrieved %v records\n", k, len(l))
	}

	if c.file != "" {
		var w bytes.Buffer
		if err := tsv(list, ctx.devices, &w); err != nil {
			return err
		}

		return os.WriteFile(c.file, w.Bytes(), 0660)
	}

	var w strings.Builder
	if err := txt(list, ctx.devices, &w); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(w.String())
	fmt.Println()

	return nil
}

func (c *GetACL) parseArgs() error {
	flagset := flag.NewFlagSet("", flag.ExitOnError)
	withPIN := flagset.Bool("with-pin", false, "Include card keypad PIN code in retrieved ACL information")
	file := ""
	args := flag.Args()[1:]

	flagset.Parse(args)

	// ... file
	if len(flagset.Args()) > 0 {
		file = flagset.Arg(0)
		stat, err := os.Stat(file)
		if err != nil && !os.IsNotExist(err) {
			return err
		} else if err == nil && stat.Mode().IsDir() {
			return fmt.Errorf("file '%s' is a directory", file)
		} else if err == nil && !stat.Mode().IsRegular() {
			return fmt.Errorf("file '%s' is not a real file", file)
		}
	}

	c.file = file
	c.withPIN = *withPIN

	return nil
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
	fmt.Println("    --with-pin Includes the card keypad PIN code in the retrieved ACL.")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli get-acl")
	fmt.Println("    uhppote-cli --debug get-acl \"uhppote.tsv\"")
	fmt.Println("    uhppote-cli --debug --config .config get-acl --with-pin \"uhppote.tsv\"")
	fmt.Println()
}

// Returns true - configuration is not optional for this command to return valid information.
func (c *GetACL) RequiresConfig() bool {
	return true
}
