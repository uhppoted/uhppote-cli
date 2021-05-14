package commands

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-api/acl"
	"github.com/uhppoted/uhppoted-api/config"
	"io"
	"os"
	"text/template"
	"time"
)

// CompareACLCmd is an initialized CompareACL command for the main() command list
var CompareACLCmd = CompareACL{
	template: `
-----------------------------------
ACL DIFF REPORT {{ .DateTime }}
{{range $id,$value := .Diffs}}
  DEVICE {{ $id }}{{if or $value.Updated $value.Added $value.Deleted}}{{else}} OK{{end}}{{if $value.Updated}}
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
		return fmt.Errorf("compare-acl requires a valid configuration file")
	}

	devices := getDevices(&ctx)

	file, err := c.getACLFile()
	if err != nil {
		return err
	}

	tsv, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	list, warnings, err := acl.ParseTSV(bytes.NewReader(tsv), devices, false)
	if err != nil {
		return err
	}

	for _, w := range warnings {
		fmt.Printf("   ... WARNING    %v\n", w)
	}

	for k, l := range list {
		fmt.Printf("   ... %v  ACL has %v records\n", k, len(l))
	}

	current, errors := acl.GetACL(ctx.uhppote, devices)
	if len(errors) > 0 {
		return fmt.Errorf("%v", errors)
	}

	diff, err := acl.Compare(current, list)
	if err != nil {
		return err
	}

	widths := map[string]int{}
	for k, v := range diff {
		if w := len(fmt.Sprintf("%v", k)); w > widths["device"] {
			widths["device"] = w
		}

		if w := len(fmt.Sprintf("%v", len(v.Unchanged))); w > widths["unchanged"] {
			widths["unchanged"] = w
		}

		if w := len(fmt.Sprintf("%v", len(v.Updated))); w > widths["updated"] {
			widths["updated"] = w
		}

		if w := len(fmt.Sprintf("%v", len(v.Added))); w > widths["added"] {
			widths["added"] = w
		}

		if w := len(fmt.Sprintf("%v", len(v.Deleted))); w > widths["deleted"] {
			widths["deleted"] = w
		}
	}

	format := fmt.Sprintf("   ... %%-%vv  same:%%-%vv  different:%%-%vv  missing:%%-%vv  extraneous:%%-%vv\n",
		widths["device"],
		widths["unchanged"],
		widths["updated"],
		widths["added"],
		widths["deleted"])

	for k, v := range diff {
		fmt.Printf(format, k, len(v.Unchanged), len(v.Updated), len(v.Added), len(v.Deleted))
	}

	var w bytes.Buffer
	if err := c.report(diff, &w); err != nil {
		return err
	}

	if rptfile, err := c.getReportFile(); err != nil {
		return err
	} else if rptfile != "" {
		return os.WriteFile(rptfile, w.Bytes(), 0660)
	}

	fmt.Printf("%s\n", string(w.Bytes()))

	return nil
}

func (c *CompareACL) report(diff map[uint32]acl.Diff, w io.Writer) error {
	t, err := template.New("report").Parse(c.template)
	if err != nil {
		return err
	}

	timestamp := types.DateTime(time.Now())
	rpt := struct {
		DateTime *types.DateTime
		Diffs    map[uint32]acl.Diff
	}{
		DateTime: &timestamp,
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

func (c *CompareACL) getReportFile() (string, error) {
	if len(flag.Args()) < 3 {
		return "", nil
	}

	file := flag.Arg(2)
	stat, err := os.Stat(file)

	if err != nil {
		if !os.IsNotExist(err) {
			return "", fmt.Errorf("Cannot use to report file '%s':%v", file, err)
		}

		return file, nil
	}

	if stat.Mode().IsDir() {
		return "", fmt.Errorf("File '%s' is a directory", file)
	}

	if !stat.Mode().IsRegular() {
		return "", fmt.Errorf("File '%s' is not a real file", file)
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
	fmt.Println(" Duplicate card numbers are ignored (with a warning message)")
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
	fmt.Printf("              (defaults to %s)\n", config.DefaultConfig)
	fmt.Println("    --debug   Displays internal information for diagnosing errors")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli --debug --config .config compare-acl \"uhppote-2019-05-25.tsv\"")
	fmt.Println()
}

// Returns true - configuration is not optional for this command to return valid information.
func (c *CompareACL) RequiresConfig() bool {
	return true
}
