package commands

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-api/config"
)

var AddTaskCmd = AddTask{}

type AddTask struct {
}

func (c *AddTask) Execute(ctx Context) error {
	serialNumber, err := getSerialNumber(ctx)
	if err != nil {
		return err
	}

	task := types.Task{
		Weekdays: types.Weekdays{},
	}

	args := flag.Args()

	// ... task type
	if len(args) < 3 {
		return fmt.Errorf("missing task identifier")
	} else {
		arg := args[2]
		if regexp.MustCompile("^[0-9]+$").MatchString(arg) {
			if taskID, err := strconv.Atoi(arg); err != nil {
				return fmt.Errorf("invalid task identifier %v (%v)", arg, err)
			} else if taskID < 0 || taskID > 12 {
				return fmt.Errorf("invalid task identifier %v - valid range is [0..12]", arg)
			} else {
				task.Task = types.TaskType(taskID)
			}
		} else {
			re := regexp.MustCompile("[^a-z]+")
			clean := func(s string) string { return re.ReplaceAllString(strings.ToLower(s), "") }
			taskID := clean(arg)
			ok := false

			for _, t := range []types.TaskType{
				types.DoorControlled,
				types.DoorOpen,
				types.DoorClosed,
				types.DisableTimeProfile,
				types.EnableTimeProfile,
				types.CardNoPassword,
				types.CardInPassword,
				types.CardInOutPassword,
				types.EnableMoreCards,
				types.DisableMoreCards,
				types.TriggerOnce,
				types.DisablePushButton,
				types.EnablePushButton,
			} {
				if taskID == clean(fmt.Sprintf("%v", t)) {
					task.Task = t
					ok = true
					break
				}
			}

			if !ok {
				return fmt.Errorf("invalid task identifier '%v'", arg)
			}
		}
	}

	// ... door
	if len(args) < 4 {
		return fmt.Errorf("missing door")
	} else {
		arg := args[3]
		if regexp.MustCompile("^[0-9]$").MatchString(arg) {
			if door, err := strconv.Atoi(arg); err != nil {
				return fmt.Errorf("invalid door '%v' (%v)", arg, err)
			} else if door < 1 || door > 4 {
				return fmt.Errorf("invalid door (%v) - valid range is [1..4]", door)
			} else {
				task.Door = uint8(door)
			}
		} else {
			return fmt.Errorf("invalid door '%v'", arg)
		}
	}

	// ... from:to
	if args := flag.Args(); len(args) < 5 {
		return fmt.Errorf("missing 'from:to' dates")
	} else {
		arg := args[4]
		re := regexp.MustCompile("([0-9]{4}-[0-9]{2}-[0-9]{2}):([0-9]{4}-[0-9]{2}-[0-9]{2})")

		match := re.FindStringSubmatch(arg)
		if match == nil || len(match) != 3 {
			return fmt.Errorf("invalid 'from:to' dates (%v)", arg)
		}

		if date, err := types.DateFromString(match[1]); err != nil {
			return fmt.Errorf("%v: invalid 'from' date (%v)", match[1], err)
		} else {
			task.From = date
		}

		if date, err := types.DateFromString(match[2]); err != nil {
			return fmt.Errorf("%v: invalid 'to' date (%v)", match[1], err)
		} else {
			task.To = date
		}
	}

	// ... weekdays
	var weekdays = days{
		"Monday":    true,
		"Tuesday":   true,
		"Wednesday": true,
		"Thursday":  true,
		"Friday":    true,
		"Saturday":  true,
		"Sunday":    true,
	}

	for _, arg := range args[5:] {
		if regexp.MustCompile("^(?i:Mon|Tue|Wed|Thu|Fri|Sat|Sun).*").MatchString(arg) {
			if err := weekdays.parse(arg); err != nil {
				return err
			}
		}
	}

	task.Weekdays = types.Weekdays{
		time.Monday:    weekdays["Monday"],
		time.Tuesday:   weekdays["Tuesday"],
		time.Wednesday: weekdays["Wednesday"],
		time.Thursday:  weekdays["Thursday"],
		time.Friday:    weekdays["Friday"],
		time.Saturday:  weekdays["Saturday"],
		time.Sunday:    weekdays["Sunday"],
	}

	// ... start time
	for _, arg := range args[5:] {
		if regexp.MustCompile("^[0-9]{2}:[0-9]{2}$").MatchString(arg) {
			if hhmm, err := types.HHmmFromString(arg); err != nil {
				return fmt.Errorf("invalid start time '%v' (%v)", arg, err)
			} else if hhmm == nil {
				return fmt.Errorf("invalid start time (%v)", arg)
			} else {
				task.Start = *hhmm
			}
		}
	}

	// ... more cards
	if task.Task == types.EnableMoreCards {
		for _, arg := range args[5:] {
			if regexp.MustCompile("^[0-9]+$").MatchString(arg) {
				if cards, err := strconv.Atoi(arg); err != nil {
					return fmt.Errorf("invalid more cards '%v' (%v)", arg, err)
				} else if cards < 0 || cards > 255 {
					return fmt.Errorf("invalid 'more cards' (%v)", arg)
				} else {
					task.Cards = uint8(cards)
				}
			}
		}
	}

	// ... good to go apparently
	if ctx.uhppote != nil && ctx.debug {
		fmt.Println(" ...")
		fmt.Printf(" ... serial number: %v\n", serialNumber)
		fmt.Printf(" ... task:          %v\n", task.Task)
		fmt.Printf(" ... door:          %v\n", task.Door)
		fmt.Printf(" ... from:          %v\n", task.From)
		fmt.Printf(" ... to:            %v\n", task.To)
		fmt.Printf(" ... weekdays:      %v\n", task.Weekdays)
		fmt.Printf(" ... start time:    %v\n", task.Start)
		fmt.Printf(" ... more cards:    %v\n", task.Cards)
		fmt.Println(" ...")
	}

	if ok, err := ctx.uhppote.AddTask(serialNumber, task); err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("%v: could not add task", serialNumber)
	}

	fmt.Printf("%v: task added\n", serialNumber)

	return nil
}

func (c *AddTask) CLI() string {
	return "add-task"
}

func (c *AddTask) Description() string {
	return "Adds a task to the controller task list"
}

func (c *AddTask) Usage() string {
	return "<serial number> <task> <door> <active> <weekdays> <start>i <cards>"
}

func (c *AddTask) Help() {
	fmt.Println("Usage: uhppote-cli [options] add-task <serial-number> <task> <active> <weekdays> <start> <cards>")
	fmt.Println()
	fmt.Println(" Adds a new task to a controller's task list")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println("  task           (required) task ID or description. The task ID or description must correspond to one of")
	fmt.Println("                            the task types listed below")
	fmt.Println("  active         (required) active start and end dates formatted as YYYY-mm-dd:YYYY-mm-dd")
	fmt.Println("  weekdays       (optional) list of weekdays on which profile is enabled (defaults to all)")
	fmt.Println("  start          (optional) start time (HH:mm) for the taskthe task (defaults to 00:00 if not defined)")
	fmt.Println("  cards          (optional) number of 'more cards' permitted for the 'enable more cards' task")
	fmt.Println()
	fmt.Println("  Tasks:")
	fmt.Println("    0   control door")
	fmt.Println("    1   unlock door")
	fmt.Println("    2   lock door")
	fmt.Println("    3   disable time profile")
	fmt.Println("    4   enable time profile")
	fmt.Println("    5   enable card, no password")
	fmt.Println("    6   enable card+IN password")
	fmt.Println("    7   enable card+password")
	fmt.Println("    8   enable more cards")
	fmt.Println("    9   disable more cards")
	fmt.Println("    10  trigger once")
	fmt.Println("    11  disable pushbutton")
	fmt.Println("    12  enable pushbutton")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --config  File path for the 'conf' file containing the controller configuration")
	fmt.Printf("              (defaults to %s)\n", config.DefaultConfig)
	fmt.Println("    --debug   Displays internal information for diagnosing errors")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli add-task 9876543210 3 2021-04-01:2021-12-31 Mon,Wed,Fri 09:30")
	fmt.Println("    uhppote-cli add-task 9876543210 'door controlled' 2021-04-01:2021-12-31 Mon,Wed,Fri 09:30")
	fmt.Println("    uhppote-cli add-task 9876543210 'enable more cards' 2021-04-01:2021-12-31 Mon,Wed,Fri 09:30 27")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *AddTask) RequiresConfig() bool {
	return false
}
