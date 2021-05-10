package commands

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"text/template"
	"time"

	"github.com/uhppoted/uhppote-core/types"
)

var GetTimeProfilesCmd = GetTimeProfiles{
	template: `
-------------------------------------------
TIME PROFILES {{.DeviceID}} {{.Timestamp}}
-------------------------------------------
Profile  From       To          Mon Tue Wed Thurs Fri Sat Sun  Start1 End1   Start2 End2   Start3 End3   Linked{{range $id,$row := .Profiles}}
{{printf "%-7s" $row.ID}}  {{printf "%-10s" $row.From}} {{printf "%-10s" $row.To}}  {{$row.Monday}}   {{$row.Tuesday}}   {{$row.Wednesday}}   {{$row.Thursday}}     {{$row.Friday}}   {{$row.Saturday}}   {{$row.Sunday}}    {{printf "%-5s" $row.Start1}}  {{printf "%-5s" $row.End1}}  {{printf "%-5s" $row.Start2}}  {{printf "%-5s" $row.End2}}  {{printf "%-5s" $row.Start3}}  {{printf "%-5s" $row.End3}}  {{printf "%-3s" $row.Linked}}{{end}}
`,
}

type GetTimeProfiles struct {
	template string
}

func (c *GetTimeProfiles) Execute(ctx Context) error {
	serialNumber, err := getSerialNumber(ctx)
	if err != nil {
		return err
	}

	f := func(v bool) string {
		if v {
			return "Y"
		} else {
			return "N"
		}
	}

	g := func(v uint8) string {
		if v == 0 {
			return ""
		} else {
			return fmt.Sprintf("%v", v)
		}
	}

	recordset := []map[string]string{}
	for id := 2; id <= 254; id++ {
		if profile, err := ctx.uhppote.GetTimeProfile(serialNumber, uint8(id)); err != nil {
			return err
		} else if profile != nil {
			row := map[string]string{
				"ID":        fmt.Sprintf("%v", profile.ID),
				"From":      fmt.Sprintf("%v", profile.From),
				"To":        fmt.Sprintf("%v", profile.To),
				"Monday":    f(profile.Weekdays[time.Monday]),
				"Tuesday":   f(profile.Weekdays[time.Tuesday]),
				"Wednesday": f(profile.Weekdays[time.Wednesday]),
				"Thursday":  f(profile.Weekdays[time.Thursday]),
				"Friday":    f(profile.Weekdays[time.Friday]),
				"Saturday":  f(profile.Weekdays[time.Saturday]),
				"Sunday":    f(profile.Weekdays[time.Sunday]),
				"Linked":    g(profile.LinkedProfileID),
			}

			if segment, ok := profile.Segments[1]; ok {
				row["Start1"] = fmt.Sprintf("%v", segment.Start)
				row["End1"] = fmt.Sprintf("%v", segment.End)
			}

			if segment, ok := profile.Segments[2]; ok {
				row["Start2"] = fmt.Sprintf("%v", segment.Start)
				row["End2"] = fmt.Sprintf("%v", segment.End)
			}

			if segment, ok := profile.Segments[3]; ok {
				row["Start3"] = fmt.Sprintf("%v", segment.Start)
				row["End3"] = fmt.Sprintf("%v", segment.End)
			}

			recordset = append(recordset, row)
		}
	}

	if file, err := c.getTSVFile(); err != nil {
		return err
	} else if file != "" {
		return c.export(file, recordset)
	}

	return c.print(serialNumber, recordset)
}

func (c *GetTimeProfiles) print(serialNumber uint32, recordset []map[string]string) error {
	timestamp := types.DateTime(time.Now())

	rpt := struct {
		DeviceID  uint32
		Timestamp *types.DateTime
		Profiles  []map[string]string
	}{
		DeviceID:  serialNumber,
		Timestamp: &timestamp,
		Profiles:  recordset,
	}

	t, err := template.New("report").Parse(c.template)
	if err != nil {
		return err
	}

	return t.Execute(os.Stdout, rpt)
}

func (c *GetTimeProfiles) export(file string, recordset []map[string]string) error {
	var b bytes.Buffer

	w := csv.NewWriter(&b)
	w.Comma = '\t'

	// TSV header
	header := []string{"Profile", "From", "To", "Mon", "Tue", "Wed", "Thurs", "Fri", "Sat", "Sun", "Start1", "End1", "Start2", "End2", "Start3", "End3", "Linked"}
	if err := w.Write(header); err != nil {
		return err
	}

	for _, r := range recordset {
		row := []string{
			r["ID"],
			r["From"],
			r["To"],
			r["Monday"],
			r["Tuesday"],
			r["Wednesday"],
			r["Thursday"],
			r["Friday"],
			r["Saturday"],
			r["Sunday"],
			r["Start1"],
			r["End1"],
			r["Start2"],
			r["End2"],
			r["Start3"],
			r["End3"],
			r["Linked"],
		}

		if err := w.Write(row); err != nil {
			return err
		}
	}

	w.Flush()

	return os.WriteFile(file, b.Bytes(), 0660)
}

func (c *GetTimeProfiles) getTSVFile() (string, error) {
	if len(flag.Args()) < 3 {
		return "", nil
	}

	file := flag.Arg(2)
	stat, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return file, nil
		} else {
			return "", err
		}
	}

	if stat.Mode().IsDir() {
		return "", fmt.Errorf("File '%s' is a directory", file)
	}

	if !stat.Mode().IsRegular() {
		return "", fmt.Errorf("File '%s' is not a real file", file)
	}

	return file, nil
}

func (c *GetTimeProfiles) CLI() string {
	return "get-time-profiles"
}

func (c *GetTimeProfiles) Description() string {
	return "Retrieves all the defined time profiles from a controller"
}

func (c *GetTimeProfiles) Usage() string {
	return "<serial number>"
}

func (c *GetTimeProfiles) Help() {
	fmt.Println("Usage: uhppote-cli [options] get-time-profiles <serial number> <file>")
	fmt.Println()
	fmt.Println(" Retrieves all the defined time profiles from a controller and (optionally) writes them to a TSV file")
	fmt.Println()
	fmt.Println("  serial number  (required) controller serial number")
	fmt.Println("  file           (optional) TSV file for time profiles")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli get-time-profiles 9876543210 ./9876543210.tsv")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *GetTimeProfiles) RequiresConfig() bool {
	return false
}
