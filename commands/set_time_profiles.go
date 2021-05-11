package commands

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/uhppoted/uhppote-cli/encoding/tsv"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-api/config"
)

var SetTimeProfilesCmd = SetTimeProfiles{}

type SetTimeProfiles struct {
}

func (c *SetTimeProfiles) Execute(ctx Context) error {
	serialNumber, err := getSerialNumber(ctx)
	if err != nil {
		return err
	}

	file, err := c.getTSVFile()
	if err != nil {
		return err
	} else if file == "" {
		return fmt.Errorf("Missing TSV file with time profiles")
	}

	profiles, err := c.parse(file)
	if err != nil {
		return err
	} else if profiles == nil {
		return fmt.Errorf("Could not parse TSV File '%s'", file)
	} else if len(profiles) == 0 {
		return fmt.Errorf("File '%s' does not contain any valid time profiles", file)
	}

	println(serialNumber)
	for _, p := range profiles {
		fmt.Printf("%v\n", p)
	}

	return nil
}

func (c *SetTimeProfiles) CLI() string {
	return "set-time-profiles"
}

func (c *SetTimeProfiles) Description() string {
	return "Writes the time profiles defined in a TSV file to a controller"
}

func (c *SetTimeProfiles) Usage() string {
	return "<serial number>"
}

func (c *SetTimeProfiles) Help() {
	fmt.Println("Usage: uhppote-cli [options] set-time-profiles <serial number> <file>")
	fmt.Println()
	fmt.Println(" Writes the time profiles defined in a TSV file to a controller. Existing time profiles are not cleared")
	fmt.Println(" but will be overwritten if redefined in the TSV file.")
	fmt.Println()
	fmt.Println("  serial number  (required) controller serial number")
	fmt.Println("  file           (required) TSV file with time profiles")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --config  File path for the 'conf' file containing the controller configuration")
	fmt.Printf("              (defaults to %s)\n", config.DefaultConfig)
	fmt.Println("    --debug   Displays internal information for diagnosing errors")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli set-time-profiles 9876543210 9876543210.tsv")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *SetTimeProfiles) RequiresConfig() bool {
	return false
}

func (c *SetTimeProfiles) getTSVFile() (string, error) {
	if len(flag.Args()) < 3 {
		return "", nil
	}

	file := flag.Arg(2)
	stat, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return file, fmt.Errorf("File '%s' does not exist", file)
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

func (c *SetTimeProfiles) parse(file string) ([]types.TimeProfile, error) {
	type tsvProfile struct {
		ID        int         `tsv:"Profile"`
		From      types.Date  `tsv:"From"`
		To        types.Date  `tsv:"To"`
		Monday    bool        `tsv:"Mon"`
		Tuesday   bool        `tsv:"Tue"`
		Wednesday bool        `tsv:"Wed"`
		Thursday  bool        `tsv:"Thurs"`
		Friday    bool        `tsv:"Fri"`
		Saturday  bool        `tsv:"Sat"`
		Sunday    bool        `tsv:"Sun"`
		Start1    *types.HHmm `tsv:"Start1"`
		End1      *types.HHmm `tsv:"End1"`
		Start2    *types.HHmm `tsv:"Start2"`
		End2      *types.HHmm `tsv:"End2"`
		Start3    *types.HHmm `tsv:"Start3"`
		End3      *types.HHmm `tsv:"End3"`
		Linked    int         `tsv:"Linked"`
	}

	recordset := []tsvProfile{}

	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if err := tsv.Unmarshal(bytes, &recordset); err != nil {
		return nil, err
	}

	profiles := []types.TimeProfile{}
	for _, record := range recordset {
		from := record.From
		to := record.To

		profile := types.TimeProfile{
			ID:              uint8(record.ID),
			LinkedProfileID: uint8(record.Linked),
			From:            &from,
			To:              &to,
			Weekdays: types.Weekdays{
				time.Monday:    record.Monday,
				time.Tuesday:   record.Tuesday,
				time.Wednesday: record.Wednesday,
				time.Thursday:  record.Thursday,
				time.Friday:    record.Friday,
				time.Saturday:  record.Saturday,
				time.Sunday:    record.Sunday,
			},
			Segments: types.Segments{
				1: types.Segment{
					Start: record.Start1,
					End:   record.End1,
				},
				2: types.Segment{
					Start: record.Start2,
					End:   record.End2,
				},
				3: types.Segment{
					Start: record.Start3,
					End:   record.End3,
				},
			},
		}

		profiles = append(profiles, profile)
	}

	return profiles, nil
}
