package commands

import (
	"flag"
	"fmt"
	"os"
	"reflect"
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
		return fmt.Errorf("Could not extract time profiles from TSV File '%s'", file)
	} else if len(profiles) == 0 {
		return fmt.Errorf("File '%s' does not contain any valid time profiles", file)
	}

	warnings, err := c.load(ctx, serialNumber, profiles)
	if err != nil {
		return err
	}

	if len(warnings) > 0 {
		fmt.Println()
		for _, warning := range warnings {
			fmt.Printf("   WARN  %v\n", warning)
		}
		fmt.Println()
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

func (c *SetTimeProfiles) load(ctx Context, serialNumber uint32, profiles []types.TimeProfile) ([]error, error) {
	prewarn := []error{}

	// check for duplicate profiles
	set := map[uint8]int{}

	for i, profile := range profiles {
		if line, ok := set[profile.ID]; ok {
			if !reflect.DeepEqual(profile, profiles[line-1]) {
				return prewarn, fmt.Errorf("Profile %v has more than one definition (records %v and %v)", profile.ID, line, i+1)
			}

			prewarn = append(prewarn, fmt.Errorf("Profile %-3v is defined twice (records %v and %v)", profile.ID, line, i+1))
		}

		set[profile.ID] = i + 1
	}

	// loop until all profiles are either set or could not be set
	warnings := prewarn[:]
	remaining := map[uint8]struct{}{}
	for _, p := range profiles {
		remaining[p.ID] = struct{}{}
	}

	for len(remaining) > 0 {
		warnings = prewarn[:]
		count := 0

		for _, profile := range profiles {
			// already loaded?
			if _, ok := remaining[profile.ID]; !ok {
				continue
			}

			// profile ok?
			if err := c.validate(profile); err != nil {
				warnings = append(warnings, fmt.Errorf("profile %-3v %v", profile.ID, err))
				continue
			}

			// verify linked profile exists
			if linked := profile.LinkedProfileID; linked != 0 {
				if p, err := ctx.uhppote.GetTimeProfile(serialNumber, linked); err != nil {
					return nil, err
				} else if p == nil {
					warnings = append(warnings, fmt.Errorf("profile %-3v linked time profile %v is not defined", profile.ID, linked))
					continue
				}
			}

			// check for circular references
			if err := c.circular(ctx, serialNumber, profile); err != nil {
				warnings = append(warnings, fmt.Errorf("profile %-3v %v", profile.ID, err))
				continue
			}

			// good to go!
			if ok, err := ctx.uhppote.SetTimeProfile(serialNumber, profile); err != nil {
				return nil, err
			} else if !ok {
				warnings = append(warnings, fmt.Errorf("%v: could not create time profile %v", serialNumber, profile.ID))
			} else {
				fmt.Printf("   ... set time profile %v\n", profile.ID)

				delete(remaining, profile.ID)
				count++
			}
		}

		if count == 0 {
			break
		}
	}

	return warnings, nil
}

func (c *SetTimeProfiles) validate(profile types.TimeProfile) error {
	if profile.From == nil {
		return fmt.Errorf("invalid 'From' date (%v)", profile.From)
	}

	if profile.To == nil {
		return fmt.Errorf("invalid 'To' date (%v)", profile.To)
	}

	if profile.To.Before(*profile.From) {
		return fmt.Errorf("'To' date (%v) is before 'From' date (%v)", profile.To, profile.From)
	}

	for _, i := range []uint8{1, 2, 3} {
		segment := profile.Segments[i]

		if segment.End.Before(segment.Start) {
			return fmt.Errorf("segment %v 'End' (%v) is before 'Start' (%v)", i, segment.End, segment.Start)
		}
	}

	return nil
}

func (c *SetTimeProfiles) circular(ctx Context, serialNumber uint32, profile types.TimeProfile) error {
	if linked := profile.LinkedProfileID; linked != 0 {
		profiles := map[uint8]bool{profile.ID: true}
		chain := []uint8{profile.ID}

		for l := linked; l != 0; {
			if p, err := ctx.uhppote.GetTimeProfile(serialNumber, l); err != nil {
				return err
			} else if p == nil {
				return fmt.Errorf("linked time profile %v is not defined", l)
			} else {
				chain = append(chain, p.ID)
				if profiles[p.ID] {
					return fmt.Errorf("linking to time profile %v creates a circular reference %v", profile.LinkedProfileID, chain)
				}

				profiles[p.ID] = true
				l = p.LinkedProfileID
			}
		}
	}

	return nil
}

func (c *SetTimeProfiles) parse(file string) ([]types.TimeProfile, error) {
	type tsvProfile struct {
		ID        int        `tsv:"Profile"`
		From      types.Date `tsv:"From"`
		To        types.Date `tsv:"To"`
		Monday    bool       `tsv:"Mon"`
		Tuesday   bool       `tsv:"Tue"`
		Wednesday bool       `tsv:"Wed"`
		Thursday  bool       `tsv:"Thurs"`
		Friday    bool       `tsv:"Fri"`
		Saturday  bool       `tsv:"Sat"`
		Sunday    bool       `tsv:"Sun"`
		Start1    types.HHmm `tsv:"Start1"`
		End1      types.HHmm `tsv:"End1"`
		Start2    types.HHmm `tsv:"Start2"`
		End2      types.HHmm `tsv:"End2"`
		Start3    types.HHmm `tsv:"Start3"`
		End3      types.HHmm `tsv:"End3"`
		Linked    int        `tsv:"Linked"`
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
