package commands

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/uhppoted/uhppote-core/uhppote"
	"github.com/uhppoted/uhppoted-lib/acl"
	"github.com/uhppoted/uhppoted-lib/config"
)

var LoadACLCmd = LoadACL{
	file:    "",
	withPIN: false,
	strict:  false,
	dryrun:  false,
}

type LoadACL struct {
	file    string
	withPIN bool
	strict  bool
	dryrun  bool
}

func (c *LoadACL) Execute(ctx Context) error {
	if ctx.config == nil {
		return errors.New("load-acl requires a valid configuration file")
	}

	if err := c.parseArgs(); err != nil {
		return err
	}

	if c.file == "" {
		return fmt.Errorf("please specify the TSV file from which to load the access control list ")
	}

	tsv, err := os.ReadFile(c.file)
	if err != nil {
		return err
	}

	list, warnings, err := acl.ParseTSV(bytes.NewReader(tsv), ctx.devices, c.strict)
	if err != nil {
		return err
	}

	for _, w := range warnings {
		fmt.Printf("   ... WARNING    %v\n", w)
	}

	for k, l := range list {
		fmt.Printf("   ... %v  ACL has %v records\n", k, len(l))
	}

	put := func(u uhppote.IUHPPOTE, list acl.ACL) (map[uint32]acl.Report, []error) {
		if c.withPIN {
			return acl.PutACLWithPIN(ctx.uhppote, list, c.dryrun)
		} else {
			return acl.PutACL(ctx.uhppote, list, c.dryrun)
		}
	}

	rpt, errors := put(ctx.uhppote, list)
	for k, v := range rpt {
		fmt.Printf("   ... %v  unchanged:%v  updated:%v  added:%v  deleted:%v  failed:%v  errors:%v\n",
			k,
			len(v.Unchanged),
			len(v.Updated),
			len(v.Added),
			len(v.Deleted),
			len(v.Failed),
			len(v.Errors))
	}

	if len(errors) > 0 {
		return fmt.Errorf("%v", errors)
	}

	fmt.Println()

loop:
	for _, r := range rpt {
		if r.Errors != nil {
			count := 0
			for _, v := range rpt {
				for _, err := range v.Errors {
					if count < 5 {
						fmt.Printf("   WARNING: %v\n", err)
						count += 1
					} else {
						fmt.Printf("   WARNING: ... etc\n")
						break loop
					}
				}
			}

			fmt.Println()
			break
		}
	}

	return nil
}

func (c *LoadACL) parseArgs() error {
	flagset := flag.NewFlagSet("", flag.ExitOnError)
	withPIN := flagset.Bool("with-pin", false, "Include card keypad PIN code in retrieved ACL information")
	strict := flagset.Bool("strict", false, "Treat duplicate card numbers as errors")
	file := ""
	args := flag.Args()[1:]

	flagset.Parse(args)

	// ... file
	if len(flagset.Args()) > 0 {
		file = flagset.Arg(0)
		stat, err := os.Stat(file)
		if err != nil && os.IsNotExist(err) {
			return fmt.Errorf("file '%s' does not exist", file)
		} else if err != nil {
			return err
		} else if err == nil && stat.Mode().IsDir() {
			return fmt.Errorf("file '%s' is a directory", file)
		} else if err == nil && !stat.Mode().IsRegular() {
			return fmt.Errorf("file '%s' is not a real file", file)
		}
	}

	c.file = file
	c.withPIN = *withPIN
	c.strict = *strict

	return nil
}

func (c *LoadACL) CLI() string {
	return "load-acl"
}

func (c *LoadACL) Description() string {
	return "Downloads an access control list from a TSV file to a set of access controllers"
}

func (c *LoadACL) Usage() string {
	return "<TSV file>"
}

func (c *LoadACL) Help() {
	fmt.Println("Usage: uhppote-cli [options] load-acl [--with-pin] [--strict] <TSV file>")
	fmt.Println()
	fmt.Println(" Downloads the access control list in the TSV file to the access controllers defined in the configuration")
	fmt.Println(" file. Duplicate card numbers are ignored (or deleted if they exist) with a warning message unless the")
	fmt.Println(" --strict option is specified")
	fmt.Println()
	fmt.Println("  <TSV file>  (required) TSV file with access control list")
	fmt.Println()
	fmt.Println("              The TSV file should conform to the following format:")
	fmt.Println("              Card Number<tab>From<tab>To<tab>Front Door<tab>Back Door<tab> ...")
	fmt.Println("              123456789<tab>2023-01-01<tab>2023-12-31<tab>Y<tab>N<tab> ...")
	fmt.Println("              987654321<tab>2023-03-05<tab>2023-11-15<tab>N<tab>N<tab> ...")
	fmt.Println()
	fmt.Println("              'Front Door', 'Back Door', etc should match the door labels in the configuration file.")
	fmt.Println("              The CLI will load the access control permissions across all the controllers listed,")
	fmt.Println("              adding cards where necessary and deleting cards not listed in the TSV file. Making")
	fmt.Println("              a backup copy of the existing permissions (using e.g. get-cards) before executing this")
	fmt.Println("              is highly recommended.")
	fmt.Println()
	fmt.Println("    --strict   Fails the load with an error if the provided ACL contains duplicate cards")
	fmt.Println("    --with-pin Updates the card keypad PIN code on the access controllers. Defaults to false.")
	fmt.Println()
	fmt.Println("               The TSV file with PIN should conform to the following format:")
	fmt.Println("               Card Number<tab>PIN<tab>From<tab>To<tab>Front Door<tab>Back Door<tab> ...")
	fmt.Println("               123456789<tab>0<tab>2023-01-01<tab>2023-12-31<tab>Y<tab>N<tab> ...")
	fmt.Println("               987654321<tab>7531<tab>2023-03-05<tab>2023-11-15<tab>N<tab>N<tab> ...")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --config  File path for the 'conf' file containing the controller configuration")
	fmt.Printf("              (defaults to %s)\n", config.DefaultConfig)
	fmt.Println("    --debug   Displays internal information for diagnosing errors")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli load-acl \"uhppote-2023-03-07.tsv\"")
	fmt.Println("    uhppote-cli --debug --config .config load-acl --with-pin \"uhppote-2023-03-07.tsv\"")
	fmt.Println()
}

// Returns true - configuration is not optional for this command to return valid information.
func (c *LoadACL) RequiresConfig() bool {
	return true
}
