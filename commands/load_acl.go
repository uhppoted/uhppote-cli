package commands

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/uhppoted/uhppoted-api/acl"
	"github.com/uhppoted/uhppoted-api/config"
	"io/ioutil"
	"os"
)

var LoadACLCmd = LoadACL{}

type LoadACL struct {
}

func (c *LoadACL) Execute(ctx Context) error {
	if ctx.config == nil {
		return errors.New("load-acl requires a valid configuration file")
	}

	devices := getDevices(&ctx)

	file, err := c.getACLFile()
	if err != nil {
		return err
	}

	tsv, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	list, err := acl.ParseTSV(bytes.NewReader(tsv), devices)
	if err != nil {
		return err
	}

	for k, l := range list {
		fmt.Printf("   ... %v  ACL has %v records\n", k, len(l))
	}

	rpt, err := acl.PutACL(ctx.uhppote, list)
	for k, v := range rpt {
		fmt.Printf("   ... %v  unchanged:%v  updated:%v  added:%v  deleted:%v  failed:%v\n", k, v.Unchanged, v.Updated, v.Added, v.Deleted, v.Failed)
	}

	return err
}

func (c *LoadACL) getACLFile() (string, error) {
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
	fmt.Println("Usage: uhppote-cli [options] load-acl <TSV file>")
	fmt.Println()
	fmt.Println(" Downloads the access control list in the TSV file to the access controllers defined in the configuration file")
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
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --config  File path for the 'conf' file containing the controller configuration")
	fmt.Printf("              (defaults to %s)\n", config.DefaultConfig)
	fmt.Println("    --debug   Displays vaguely useful internal information")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli --debug --config .config load-acl \"uhppote-2019-05-25.tsv\"")
	fmt.Println()
}
