package commands

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/uhppoted/uhppoted-api/acl"
	"github.com/uhppoted/uhppoted-api/config"
	"os"
)

var LoadACLCmd = LoadACL{
	strict: false,
}

type LoadACL struct {
	strict bool
}

func (c *LoadACL) Execute(ctx Context) error {
	if ctx.config == nil {
		return errors.New("load-acl requires a valid configuration file")
	}

	devices := getDevices(&ctx)

	file, err := c.parseArgs()
	if err != nil {
		return err
	}

	tsv, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	list, warnings, err := acl.ParseTSV(bytes.NewReader(tsv), devices, c.strict)
	if err != nil {
		return err
	}

	for _, w := range warnings {
		fmt.Printf("   ... WARNING    %v\n", w)
	}

	for k, l := range list {
		fmt.Printf("   ... %v  ACL has %v records\n", k, len(l))
	}

	rpt, errors := acl.PutACL(ctx.uhppote, list, false)
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

func (c *LoadACL) parseArgs() (string, error) {
	if len(flag.Args()) < 2 {
		return "", fmt.Errorf("Please specify the TSV file from which to load the access control list ")
	}

	file := flag.Arg(1)
	stat, err := os.Stat(file)

	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("File '%s' does not exist", file)
		}

		return "", fmt.Errorf("Failed to find file '%s':%v", file, err)
	}

	if stat.Mode().IsDir() {
		return "", fmt.Errorf("File '%s' is a directory", file)
	}

	if !stat.Mode().IsRegular() {
		return "", fmt.Errorf("File '%s' is not a real file", file)
	}

	for _, f := range flag.Args() {
		if f == "--strict" {
			c.strict = true
		}
	}
	return file, nil
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
	fmt.Println("Usage: uhppote-cli [options] load-acl <TSV file> [--strict]")
	fmt.Println()
	fmt.Println(" Downloads the access control list in the TSV file to the access controllers defined in the configuration")
	fmt.Println(" file. Duplicate card numbers are ignored (or deleted if they exist) with a warning message unless the")
	fmt.Println(" --strict option is specified")
	fmt.Println()
	fmt.Println("  <TSV file>  (required) TSV file with access control list")
	fmt.Println()
	fmt.Println("              The TSV file should conform to the following format:")
	fmt.Println("              Card Number<tab>From<tab>To<tab>Front Door<tab>Back Door<tab> ...")
	fmt.Println("              123456789<tab>2019-01-01<tab>2019-12-31<tab>Y<tab>N<tab> ...")
	fmt.Println("              987654321<tab>2019-03-05<tab>2019-11-15<tab>N<tab>N<tab> ...")
	fmt.Println()
	fmt.Println("              'Front Door', 'Back Door', etc should match the door labels in the configuration file.")
	fmt.Println("              The CLI will load the access control permissions across all the controllers listed,")
	fmt.Println("              adding cards where necessary and deleting cards not listed in the TSV file. Making")
	fmt.Println("              a backup copy of the existing permissions (using e.g. get-cards) before executing this")
	fmt.Println("              is highly recommended.")
	fmt.Println()
	fmt.Println("    --strict  Fails the load with an error if the provided ACL contains duplicate cards")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --config  File path for the 'conf' file containing the controller configuration")
	fmt.Printf("              (defaults to %s)\n", config.DefaultConfig)
	fmt.Println("    --debug   Displays internal information for diagnosing errors")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli --debug --config .config load-acl \"uhppote-2019-05-25.tsv\"")
	fmt.Println()
}

// Returns true - configuration is not optional for this command to return valid information.
func (c *LoadACL) RequiresConfig() bool {
	return true
}
