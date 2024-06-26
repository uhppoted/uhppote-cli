package commands

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-lib/config"
)

var SetTimeProfileCmd = SetTimeProfile{}

type SetTimeProfile struct {
}

func (c *SetTimeProfile) Execute(ctx Context) error {
	serialNumber, err := getSerialNumber(ctx)
	if err != nil {
		return err
	}

	profileID, err := getUint8(2, "missing time profile ID", "invalid time profile ID: %v")
	if err != nil {
		return err
	} else if profileID < 2 || profileID > 254 {
		return fmt.Errorf("invalid time profile ID (%v) - valid range is from 2 to 254", profileID)
	}

	var from types.Date
	var to types.Date

	var weekdays = days{
		"Monday":    true,
		"Tuesday":   true,
		"Wednesday": true,
		"Thursday":  true,
		"Friday":    true,
		"Saturday":  true,
		"Sunday":    true,
	}

	var schedule = segments{}
	var linked uint8

	if args := flag.Args(); len(args) > 2 {
		for _, arg := range args[3:] {
			// from:to
			if match := regexp.MustCompile("([0-9]{4}-[0-9]{2}-[0-9]{2}):([0-9]{4}-[0-9]{2}-[0-9]{2})").FindStringSubmatch(arg); match != nil {
				if date, err := types.ParseDate(match[1]); err != nil {
					return fmt.Errorf("%v: invalid 'start' date (%v)", match[1], err)
				} else {
					from = date
				}

				if date, err := types.ParseDate(match[2]); err != nil {
					return fmt.Errorf("%v: invalid 'to' date (%v)", match[1], err)
				} else {
					to = date
				}

			}

			// weekdays
			if regexp.MustCompile("^(?i:Mon|Tue|Wed|Thu|Fri|Sat|Sun).*").MatchString(arg) {
				if err := weekdays.parse(arg); err != nil {
					return err
				}
			}

			// segments
			if regexp.MustCompile("[0-9]{2}:[0-9]{2}-[0-9]{2}:[0-9]{2}").MatchString(arg) {
				if err := schedule.parse(arg); err != nil {
					return err
				}
			}

			// linked profile
			if regexp.MustCompile("^[0-9]+$").MatchString(arg) {
				if v, err := strconv.ParseUint(arg, 10, 8); err != nil {
					return fmt.Errorf("%v: invalid linked profile (%v)", arg, err)
				} else if v != 0 && v < 2 || v > 254 {
					return fmt.Errorf("%v: invalid linked profile (valid range is from 2 to 254)", arg)
				} else if uint8(v) == profileID {
					return fmt.Errorf("%v: invalid linked profile (link to self creates circular reference)", arg)
				} else {
					linked = uint8(v)
				}
			}
		}
	}

	if ctx.uhppote != nil && ctx.debug {
		fmt.Println(" ...")
		fmt.Printf(" ... serial number: %v\n", serialNumber)
		fmt.Printf(" ... profile ID:    %v\n", profileID)
		fmt.Printf(" ... from:          %v\n", from)
		fmt.Printf(" ... to:            %v\n", to)
		fmt.Printf(" ... weekdays:      %v\n", weekdays)
		fmt.Printf(" ... schedule:      %v\n", schedule)
		fmt.Printf(" ... linked:        %v\n", linked)
		fmt.Println(" ...")
	}

	if from.IsZero() {
		return fmt.Errorf("missing 'from' date")
	}

	if to.IsZero() {
		return fmt.Errorf("missing 'to' date")
	}

	if linked != 0 {
		if profile, err := ctx.uhppote.GetTimeProfile(serialNumber, linked); err != nil {
			return err
		} else if profile == nil {
			return fmt.Errorf("linked time profile %v is not defined", linked)
		}

		profiles := map[uint8]bool{profileID: true}
		links := []uint8{profileID}
		for l := linked; l != 0; {
			if profile, err := ctx.uhppote.GetTimeProfile(serialNumber, l); err != nil {
				return err
			} else if profile == nil {
				return fmt.Errorf("linked time profile %v is not defined", l)
			} else {
				links = append(links, profile.ID)
				if profiles[profile.ID] {
					return fmt.Errorf("linking to time profile %v creates a circular reference (%v)", linked, links)
				}

				profiles[profile.ID] = true
				l = profile.LinkedProfileID
			}
		}
	}

	profile := types.TimeProfile{
		ID:              profileID,
		LinkedProfileID: linked,
		From:            from,
		To:              to,

		Weekdays: types.Weekdays{
			time.Monday:    weekdays["Monday"],
			time.Tuesday:   weekdays["Tuesday"],
			time.Wednesday: weekdays["Wednesday"],
			time.Thursday:  weekdays["Thursday"],
			time.Friday:    weekdays["Friday"],
			time.Saturday:  weekdays["Saturday"],
			time.Sunday:    weekdays["Sunday"],
		},

		Segments: types.Segments{
			1: types.Segment{},
			2: types.Segment{},
			3: types.Segment{},
		},
	}

	for _, ix := range []int{1, 2, 3} {
		if s, ok := schedule[ix]; ok {
			segment := types.Segment{}

			if s.start != nil {
				segment.Start = *s.start
			}

			if s.end != nil {
				segment.End = *s.end
			}

			profile.Segments[uint8(ix)] = segment
		}
	}

	if ok, err := ctx.uhppote.SetTimeProfile(serialNumber, profile); err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("%v: could not create time profile %v", serialNumber, profileID)
	}

	fmt.Printf("%v: time profile %v created\n", serialNumber, profileID)

	return nil
}

func (c *SetTimeProfile) CLI() string {
	return "set-time-profile"
}

func (c *SetTimeProfile) Description() string {
	return "Sets the time profile associated with a time profile ID"
}

func (c *SetTimeProfile) Usage() string {
	return "<serial number> <profile ID> <active> <weekdays> <segment1> <segment2> <segment3> <linked profile>"
}

func (c *SetTimeProfile) Help() {
	fmt.Println("Usage: uhppote-cli [options] set-time-profile <serial-number> <profile-ID> <active> <weekdays> <segments> <linked>")
	fmt.Println()
	fmt.Println(" Retrieves the time profile associated with a profile ID")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println("  profile-ID     (required) time profile ID (2-254)")
	fmt.Println("  active         (required) active start and end dates formatted as YYYY-mm-dd:YYYY-mm-dd")
	fmt.Println("  weekdays       (optional) list of weekdays on which profile is enabled (defaults to all)")
	fmt.Println("  segments       (optional) start and end times (HH:mm-HH:mm) for up to 3 segments (segments default to 00:00-00:00 if not defined)")
	fmt.Println("  linked         (optional) ID of linked profile.Defaults to 0 (unlinked)")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --config  File path for the 'conf' file containing the controller configuration")
	fmt.Printf("              (defaults to %s)\n", config.DefaultConfig)
	fmt.Println("    --debug   Displays internal information for diagnosing errors")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli set-time-profile 9876543210 7 2021-04-01:2021-12-31 Mon,Wed,Fri 09:30-11:15,,15:45-17:30 27")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *SetTimeProfile) RequiresConfig() bool {
	return false
}

// SEGMENTS
type segments map[int]segment

func (ss segments) String() string {
	list := []string{}
	for _, ix := range []int{1, 2, 3} {
		if s, ok := ss[ix]; ok {
			list = append(list, s.String())
		} else {
			list = append(list, "")
		}
	}

	return strings.Join(list, ",")
}

type segment struct {
	start *types.HHmm
	end   *types.HHmm
}

func (s segment) String() string {
	return fmt.Sprintf("%v:%v", s.start, s.end)
}

func (ss segments) parse(arg string) error {
	tokens := strings.Split(arg, ",")
	for i, t := range tokens {
		if t == "" {
			continue
		}

		if match := regexp.MustCompile("([0-9]{2}:[0-9]{2})-([0-9]{2}:[0-9]{2})").FindStringSubmatch(t); match == nil {
			return fmt.Errorf("%v: invalid 'segment'", t)
		} else {
			start, err := types.HHmmFromString(match[1])
			if err != nil || start == nil {
				return fmt.Errorf("segment %v: invalid 'start' (%v:%v)", i+1, t, err)
			}

			end, err := types.HHmmFromString(match[2])
			if err != nil || end == nil {
				return fmt.Errorf("segment %v: invalid 'end' (%v)", i+1, t)
			}

			if i < 3 {
				ss[i+1] = segment{
					start: start,
					end:   end,
				}
			}
		}
	}

	return nil
}
