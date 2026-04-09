package commands

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/uhppoted/uhppote-core/types"
)

var SetFirstCardCmd = SetFirstCard{}

type SetFirstCard struct {
}

func (c *SetFirstCard) Execute(ctx Context) error {
	if serialNumber, err := getSerialNumber(ctx); err != nil {
		return err
	} else if door, firstcard, err := c.parse(); err != nil {
		return err
	} else if ok, err := ctx.uhppote.SetFirstCard(serialNumber, door, firstcard); err != nil {
		return err
	} else if ok {
		fmt.Printf("%v  %v set first-card ok\n", serialNumber, door)
	} else {
		fmt.Printf("%v  %v set first-card failed\n", serialNumber, door)
	}

	return nil
}

func (c *SetFirstCard) CLI() string {
	return "set-firstcard"
}

func (c *SetFirstCard) Description() string {
	return "Sets the controller first card configuration"
}

func (c *SetFirstCard) Usage() string {
	return "<serial number> <door> <start> <end> <active> <inactive> <weekdays>"
}

func (c *SetFirstCard) Help() {
	fmt.Println("Usage: uhppote-cli [options] set-firstcard <serial-number> <door> <start> <end> <active> <inactive> <weekdays> [command options]")
	fmt.Println()
	fmt.Println(" Sets the controller first card mode")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println("  door           (required) door ID ([1..4])")
	fmt.Println("  start          (required) time from which 'first card' is enabled (HH:mm)")
	fmt.Println("  end            (required) time after which 'first card' is disabled (HH:mm)")
	fmt.Println("  active         (required) door control mode after 'first card' swipe (controlled, normally-open, normally-closed)")
	fmt.Println("  inactive       (required) door control mode after 'first card' end time (controlled, normally-open, normally-closed, firstcard)")
	fmt.Println("  weekdays       (required) list of weekdays on which 'first card' is enabled (e.g. Mon, Tue, Fri)")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    -debug  Displays internal information for diagnosing errors")
	fmt.Println()
	fmt.Println("  Command options:")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *SetFirstCard) RequiresConfig() bool {
	return false
}

func (c SetFirstCard) parse() (uint8, types.FirstCard, error) {
	args := flag.Args()

	door := uint8(0)
	firstcard := types.FirstCard{
		StartTime: types.NewHHmm(0, 0),
		EndTime:   types.NewHHmm(0, 0),
		Active:    types.ModeUnknown,
		Inactive:  types.ModeUnknown,
		Weekdays:  types.Weekdays{},
	}

	modes := map[string]types.ControlState{
		"controlled":      types.ModeControlled,
		"normally-open":   types.ModeNormallyOpen,
		"normally-closed": types.ModeNormallyClosed,
		"firstcard":       types.ModeFirstCardOnly,
	}

	weekdays := map[string]time.Weekday{
		"mon": time.Monday,
		"tue": time.Tuesday,
		"wed": time.Wednesday,
		"thu": time.Thursday,
		"fri": time.Friday,
		"sat": time.Saturday,
		"sun": time.Sunday,
	}

	// ... door
	if len(args) <= 2 {
		return door, firstcard, fmt.Errorf("missing door ID")
	} else if v, err := strconv.ParseUint(args[2], 10, 8); err != nil {
		return door, firstcard, fmt.Errorf("invalid door ID (%v)", args[2])
	} else {
		door = uint8(v)
	}

	// ... start time
	if len(args) > 3 {
		if v, err := types.ParseHHmm(args[3]); err != nil {
			return door, firstcard, fmt.Errorf("invalid start-time (%v)", args[3])
		} else if v == nil {
			return door, firstcard, fmt.Errorf("invalid start-time (%v)", args[3])
		} else {
			firstcard.StartTime = *v
		}

		// ... end time
		if len(args) <= 4 {
			return door, firstcard, fmt.Errorf("missing end-time")
		} else if v, err := types.ParseHHmm(args[4]); err != nil {
			return door, firstcard, fmt.Errorf("invalid end-time (%v)", args[4])
		} else if v == nil {
			return door, firstcard, fmt.Errorf("invalid end-time (%v)", args[4])
		} else {
			firstcard.EndTime = *v
		}

		// ... active control state
		if len(args) <= 5 {
			return door, firstcard, fmt.Errorf("missing 'active' control state")
		} else if v, ok := modes[strings.ToLower(args[5])]; !ok || v == types.ModeFirstCardOnly {
			return door, firstcard, fmt.Errorf("invalid 'active' control state (%v)", args[5])
		} else {
			firstcard.Active = v
		}

		// ... inactive control state
		if len(args) <= 6 {
			return door, firstcard, fmt.Errorf("missing 'inactive' control state")
		} else if v, ok := modes[strings.ToLower(args[6])]; !ok {
			return door, firstcard, fmt.Errorf("invalid 'inactive' control state (%v)", args[6])
		} else {
			firstcard.Inactive = v
		}

		// ... weekdays
		if len(args) <= 7 {
			return door, firstcard, fmt.Errorf("missing weekdays")
		} else {
			tokens := strings.SplitSeq(strings.ToLower(args[7]), ",")

			for t := range tokens {
				if len(t) > 3 {
					t = t[:3]
				}

				if v, ok := weekdays[t]; ok {
					firstcard.Weekdays[v] = true
				}
			}
		}
	}

	return door, firstcard, nil
}
