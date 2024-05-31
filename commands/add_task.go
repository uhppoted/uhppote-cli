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

var AddTaskCmd = AddTask{}

type AddTask struct {
}

func (c *AddTask) Execute(ctx Context) error {
	serialNumber, err := getSerialNumber(ctx)
	if err != nil {
		return err
	}

	task, err := c.getTask(flag.Args())
	if err != nil {
		return err
	} else if task == nil {
		return fmt.Errorf("failed to parse task")
	}

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

	if ok, err := ctx.uhppote.AddTask(serialNumber, *task); err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("%v: failed to add task", serialNumber)
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
	fmt.Println("Usage: uhppote-cli [options] add-task <serial-number> <door> <task> <active> <weekdays> <start> <cards>")
	fmt.Println()
	fmt.Println(" Adds a new task to a controller's task list")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println("  task           (required) task ID or description. The task ID or description must correspond to one of")
	fmt.Println("                            the task types listed below")
	fmt.Println("  door           (required) door (1-4) to which task is assigned")
	fmt.Println("  active         (required) active start and end dates formatted as YYYY-mm-dd:YYYY-mm-dd")
	fmt.Println("  weekdays       (optional) list of weekdays on which profile is enabled (defaults to all)")
	fmt.Println("  start          (optional) start time (HH:mm) for the taskthe task (defaults to 00:00 if not defined)")
	fmt.Println("  cards          (optional) number of 'more cards' permitted for the 'enable more cards' task")
	fmt.Println()
	fmt.Println("  Tasks:")
	fmt.Println("    1   control door")
	fmt.Println("    2   unlock door")
	fmt.Println("    3   lock door")
	fmt.Println("    4   disable time profile")
	fmt.Println("    5   enable time profile")
	fmt.Println("    6   enable card, no password")
	fmt.Println("    7   enable card+IN password")
	fmt.Println("    8   enable card+password")
	fmt.Println("    9   enable more cards")
	fmt.Println("    10  disable more cards")
	fmt.Println("    11  trigger once")
	fmt.Println("    12  disable pushbutton")
	fmt.Println("    13  enable pushbutton")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --config  File path for the 'conf' file containing the controller configuration")
	fmt.Printf("              (defaults to %s)\n", config.DefaultConfig)
	fmt.Println("    --debug   Displays internal information for diagnosing errors")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli add-task 9876543210 3 1 2021-04-01:2021-12-31 Mon,Wed,Fri 09:30")
	fmt.Println("    uhppote-cli add-task 9876543210 'door controlled'   1 2021-04-01:2021-12-31 Mon,Wed,Fri 09:30")
	fmt.Println("    uhppote-cli add-task 9876543210 'enable more cards' 1 2021-04-01:2021-12-31 Mon,Wed,Fri 09:30 27")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *AddTask) RequiresConfig() bool {
	return false
}

func (c *AddTask) getTask(args []string) (*types.Task, error) {
	task := types.Task{
		Weekdays: types.Weekdays{},
	}

	if t, err := c.getTaskType(args); err != nil {
		return nil, err
	} else {
		task.Task = t
	}

	if d, err := c.getTaskDoor(args); err != nil {
		return nil, err
	} else {
		task.Door = d
	}

	if from, to, err := c.getTaskActive(args); err != nil {
		return nil, err
	} else if from == nil {
		return nil, fmt.Errorf("invalid 'from' date")
	} else if to == nil {
		return nil, fmt.Errorf("invalid 'to' date")
	} else {
		task.From = *from
		task.To = *to
	}

	if weekdays, err := c.getTaskDays(args); err != nil {
		return nil, err
	} else if weekdays == nil {
		return nil, fmt.Errorf("invalid list of weekdays")
	} else {
		task.Weekdays = *weekdays
	}

	if hhmm, err := c.getTaskStart(args); err != nil {
		return nil, err
	} else if hhmm != nil {
		task.Start = *hhmm
	}

	if task.Task == types.EnableMoreCards {
		if cards, err := c.getTaskCards(args); err != nil {
			return nil, err
		} else {
			task.Cards = cards
		}
	}

	return &task, nil
}

func (c *AddTask) getTaskType(args []string) (types.TaskType, error) {
	if len(args) < 3 {
		return 0, fmt.Errorf("missing task identifier")
	}

	arg := args[2]

	// ... numeric task type?
	if regexp.MustCompile("^[0-9]+$").MatchString(arg) {
		taskID, err := strconv.Atoi(arg)
		if err != nil {
			return 0, fmt.Errorf("invalid task identifier %v (%v)", arg, err)
		}

		if taskID < 1 || taskID > 13 {
			return 0, fmt.Errorf("invalid task identifier %v - valid range is [0..12]", arg)
		}

		return types.TaskType(taskID - 1), nil
	}

	// ... text task type
	re := regexp.MustCompile("[^a-z]+")
	clean := func(s string) string { return re.ReplaceAllString(strings.ToLower(s), "") }
	taskID := clean(arg)

	for _, t := range []types.TaskType{
		types.DoorControlled,
		types.DoorNormallyOpen,
		types.DoorNormallyClosed,
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
			return t, nil
		}
	}

	return 0, fmt.Errorf("invalid task identifier '%v'", arg)
}

func (c *AddTask) getTaskDoor(args []string) (uint8, error) {
	if len(args) < 4 {
		return 0, fmt.Errorf("missing door")
	}

	arg := args[3]

	if regexp.MustCompile("^[1-4]$").MatchString(arg) {
		if door, err := strconv.Atoi(arg); err != nil {
			return 0, fmt.Errorf("invalid door '%v' (%v)", arg, err)
		} else if door < 1 || door > 4 {
			return 0, fmt.Errorf("invalid door (%v) - valid range is [1..4]", door)
		} else {
			return uint8(door), nil
		}
	}

	return 0, fmt.Errorf("invalid door '%v'", arg)
}

func (c *AddTask) getTaskActive(args []string) (*types.Date, *types.Date, error) {
	if args := flag.Args(); len(args) < 5 {
		return nil, nil, fmt.Errorf("missing 'from:to' dates")
	}

	arg := args[4]
	match := regexp.MustCompile("([0-9]{4}-[0-9]{2}-[0-9]{2}):([0-9]{4}-[0-9]{2}-[0-9]{2})").FindStringSubmatch(arg)
	if match == nil || len(match) != 3 {
		return nil, nil, fmt.Errorf("invalid 'from:to' dates (%v)", arg)
	}

	var from *types.Date
	var to *types.Date

	if date, err := types.ParseDate(match[1]); err != nil {
		return nil, nil, fmt.Errorf("%v: invalid 'from' date (%v)", match[1], err)
	} else {
		from = &date
	}

	if date, err := types.ParseDate(match[2]); err != nil {
		return nil, nil, fmt.Errorf("%v: invalid 'to' date (%v)", match[1], err)
	} else {
		to = &date
	}

	return from, to, nil
}

func (c *AddTask) getTaskDays(args []string) (*types.Weekdays, error) {
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
				return nil, err
			}
		}
	}

	return &types.Weekdays{
		time.Monday:    weekdays["Monday"],
		time.Tuesday:   weekdays["Tuesday"],
		time.Wednesday: weekdays["Wednesday"],
		time.Thursday:  weekdays["Thursday"],
		time.Friday:    weekdays["Friday"],
		time.Saturday:  weekdays["Saturday"],
		time.Sunday:    weekdays["Sunday"],
	}, nil
}

func (c *AddTask) getTaskStart(args []string) (*types.HHmm, error) {
	for _, arg := range args[5:] {
		if regexp.MustCompile("^[0-9]{2}:[0-9]{2}$").MatchString(arg) {
			if hhmm, err := types.HHmmFromString(arg); err != nil {
				return nil, fmt.Errorf("invalid start time '%v' (%v)", arg, err)
			} else if hhmm == nil {
				return nil, fmt.Errorf("invalid start time (%v)", arg)
			} else {
				return hhmm, nil
			}
		}
	}
	return &types.HHmm{}, nil
}

func (c *AddTask) getTaskCards(args []string) (uint8, error) {
	for _, arg := range args[5:] {
		if regexp.MustCompile("^[0-9]+$").MatchString(arg) {
			if cards, err := strconv.Atoi(arg); err != nil {
				return 0, fmt.Errorf("invalid more cards '%v' (%v)", arg, err)
			} else if cards < 0 || cards > 255 {
				return 0, fmt.Errorf("invalid 'more cards' (%v)", arg)
			} else {
				return uint8(cards), nil
			}
		}
	}

	return 0, nil
}
