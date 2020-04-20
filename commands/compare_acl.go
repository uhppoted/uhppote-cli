package commands

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-api/acl"
	"io"
	"io/ioutil"
	"os"
	"text/template"
	"time"
)

var COMPARE_ACL = CompareACL{
	template: `
-----------------------------------
ACL DIFF REPORT {{ .DateTime }}
{{range $id,$value := .Diffs}}
  DEVICE {{ $id }}{{if $value.Unchanged}}
    Incorrect:  {{range $value.Updated}}{{.}}
                {{end}}{{end}}{{if $value.Added}}
    Missing:    {{range $value.Added}}{{.}}
                {{end}}{{end}}{{if $value.Deleted}}
    Unexpected: {{range $value.Deleted}}{{.}}
                {{end}}{{end}}{{end}}
-----------------------------------
`,
}

type CompareACL struct {
	template string
}

func (c *CompareACL) Execute(ctx Context) error {
	if ctx.config == nil {
		return errors.New("compare-acl requires a valid configuration file")
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

	current, err := acl.GetACL(ctx.uhppote, devices)
	if err != nil {
		return err
	}

	diff, err := acl.Compare(current, list)
	if err != nil {
		return err
	}

	for k, v := range diff {
		fmt.Printf("   ... %v  same:%v  different:%v  missing:%v  extraneous:%v\n", k, len(v.Unchanged), len(v.Updated), len(v.Added), len(v.Deleted))
	}

	var w bytes.Buffer
	if err := c.report(diff, &w); err != nil {
		return err
	}

	if rptfile, err := c.getReportFile(); err != nil {
		return err
	} else if rptfile != "" {
		return ioutil.WriteFile(rptfile, w.Bytes(), 0660)
	}

	fmt.Printf("%s\n", string(w.Bytes()))
	return nil
}

func (c *CompareACL) report(diff map[uint32]acl.Diff, w io.Writer) error {
	t, err := template.New("report").Parse(c.template)
	if err != nil {
		return err
	}

	rpt := struct {
		DateTime types.DateTime
		Diffs    map[uint32]acl.Diff
	}{
		DateTime: types.DateTime(time.Now()),
		Diffs:    diff,
	}

	return t.Execute(w, rpt)
}

func (c *CompareACL) getACLFile() (string, error) {
	if len(flag.Args()) < 2 {
		return "", fmt.Errorf("Please specify the TSV file from which to load the authoritative access control list ")
	}

	file := flag.Arg(1)
	stat, err := os.Stat(file)

	if err != nil {
		if os.IsNotExist(err) {
			return "", errors.New(fmt.Sprintf("File '%s' does not exist", file))
		}

		return "", errors.New(fmt.Sprintf("Failed to find file '%s':%v", file, err))
	}

	if stat.Mode().IsDir() {
		return "", errors.New(fmt.Sprintf("File '%s' is a directory", file))
	}

	if !stat.Mode().IsRegular() {
		return "", errors.New(fmt.Sprintf("File '%s' is not a real file", file))
	}

	return file, nil
}

func (c *CompareACL) getReportFile() (string, error) {
	if len(flag.Args()) < 3 {
		return "", nil
	}

	file := flag.Arg(2)
	stat, err := os.Stat(file)

	if err != nil {
		if !os.IsNotExist(err) {
			return "", errors.New(fmt.Sprintf("Cannot use to report file '%s':%v", file, err))
		}

		return file, nil
	}

	if stat.Mode().IsDir() {
		return "", errors.New(fmt.Sprintf("File '%s' is a directory", file))
	}

	if !stat.Mode().IsRegular() {
		return "", errors.New(fmt.Sprintf("File '%s' is not a real file", file))
	}

	return file, nil
}

func (c *CompareACL) CLI() string {
	return "compare-acl"
}

func (c *CompareACL) Description() string {
	return "Compares the card lists in the configured controllers to an authoritative access control list from a TSV file"
}

func (c *CompareACL) Usage() string {
	return "<TSV file> [<report file>]"
}

func (c *CompareACL) Help() {
	fmt.Println("Usage: uhppote-cli [options] compare-acl <TSV file> <report file>")
	fmt.Println()
	fmt.Println(" Compares the card lists in the configurated controllers to the authoritative access control list in the TSV file")
	fmt.Println()
	fmt.Println("  <TSV file>    (required) TSV file with access control list")
	fmt.Println()
	fmt.Println("                The TSV file should conform to the following format:")
	fmt.Println("                Card Number<tab>From<tab>To<tab>Front Door<tab>Back Door<tab> ...")
	fmt.Println("                123456789<tab>2019-01-01<tab>2019-12-31<tab>Y<tab>N<tab> ...")
	fmt.Println("                987654321<tab>2019-03-05<tab>2019-11-15<tab>N<tab>N<tab> ...")
	fmt.Println()
	fmt.Println("                'Front Door', 'Back Door', etc should match the door labels in the configuration file.")
	fmt.Println("                The CLI will compare the access control permissions across all the controllers listed.")
	fmt.Println()
	fmt.Println("  <report file> (optional) file to which to write the 'compare' report. Defaults to stdout if not provided")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --config  File path for the 'conf' file containing the controller configuration")
	fmt.Printf("              (defaults to %s)\n", DEFAULT_CONFIG)
	fmt.Println("    --debug   Displays vaguely useful internal information")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli --debug --config .config compare-acl \"uhppote-2019-05-25.tsv\"")
	fmt.Println()
}
