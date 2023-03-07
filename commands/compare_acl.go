package commands

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-lib/acl"
	"github.com/uhppoted/uhppoted-lib/config"
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

	file, rptfile, withPIN, err := c.parseArgs()
	if err != nil {
		return err
	}

	if file == "" {
		return fmt.Errorf("please specify the TSV file from which to load the authoritative access control list ")
	}

	tsv, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	list, warnings, err := acl.ParseTSV(bytes.NewReader(tsv), ctx.devices, false)
	if err != nil {
		return err
	}

	for _, w := range warnings {
		fmt.Printf("   ... WARNING    %v\n", w)
	}

	for k, l := range list {
		fmt.Printf("   ... %v  ACL has %v records\n", k, len(l))
	}

	current, errors := acl.GetACL(ctx.uhppote, ctx.devices)
	if len(errors) > 0 {
		return fmt.Errorf("%v", errors)
	}

	compare := func(current acl.ACL, list acl.ACL) (map[uint32]acl.Diff, error) {
		if withPIN {
			return acl.CompareWithPIN(current, list)
		} else {
			return acl.Compare(current, list)
		}
	}

	diff, err := compare(current, list)
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

	if rptfile != "" {
		return os.WriteFile(rptfile, w.Bytes(), 0660)
	}

	fmt.Printf("%v\n", w.String())

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

func (c *CompareACL) parseArgs() (string, string, bool, error) {
	flagset := flag.NewFlagSet("", flag.ExitOnError)
	withPIN := flagset.Bool("with-pin", false, "Include card keypad PIN code in retrieved ACL information")
	file := ""
	rptfile := ""
	args := flag.Args()[1:]

	flagset.Parse(args)

	// ... file
	if len(flagset.Args()) > 0 {
		file = flagset.Arg(0)
		stat, err := os.Stat(file)
		if err != nil && os.IsNotExist(err) {
			return "", "", false, fmt.Errorf("file '%s' does not exist", file)
		} else if err != nil {
			return "", "", false, err
		} else if err == nil && stat.Mode().IsDir() {
			return "", "", false, fmt.Errorf("file '%s' is a directory", file)
		} else if err == nil && !stat.Mode().IsRegular() {
			return "", "", false, fmt.Errorf("file '%s' is not a real file", file)
		}
	}

	// ... report file
	if len(flagset.Args()) > 1 {
		rptfile = flagset.Arg(1)
		stat, err := os.Stat(rptfile)
		if err != nil && !os.IsNotExist(err) {
			return "", "", false, err
		} else if err == nil && stat.Mode().IsDir() {
			return "", "", false, fmt.Errorf("file '%s' is a directory", rptfile)
		} else if err == nil && !stat.Mode().IsRegular() {
			return "", "", false, fmt.Errorf("file '%s' is not a real file", rptfile)
		}
	}

	return file, rptfile, *withPIN, nil
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
	fmt.Println("    --with-pin Includes the card keypad PIN code when comparing ACLs.")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli compare-acl \"uhppote-2023-03-07.tsv\"")
	fmt.Println("    uhppote-cli --debug --config .config compare-acl --with-pin \"uhppote-2023-03-07.tsv\"")
	fmt.Println()
}

// Returns true - configuration is not optional for this command to return valid information.
func (c *CompareACL) RequiresConfig() bool {
	return true
}
