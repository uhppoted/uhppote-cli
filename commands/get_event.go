package commands

import (
	"flag"
	"fmt"
)

var GetEventCmd = GetEvent{}

type GetEvent struct {
}

func (c *GetEvent) Execute(ctx Context) error {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return err
	}

	index, err := c.getIndex(ctx, serialNumber)
	if err != nil {
		return err
	}

	if index > 100000 {
		index = 1
	}

	event, err := ctx.uhppote.GetEvent(serialNumber, index)
	if err != nil {
		return err
	}

	if event == nil {
		return fmt.Errorf("%v:  no event at index: %v", serialNumber, index)
	}

	if event.Index != index {
		return fmt.Errorf("%v:  event index %v out of range", serialNumber, index)
	}

	if len(flag.Args()) < 3 {
		_, err := ctx.uhppote.SetEventIndex(serialNumber, index)
		if err != nil {
			return err
		}
	}

	fmt.Printf("%v\n", event)

	return nil
}

func (c *GetEvent) CLI() string {
	return "get-event"
}

func (c *GetEvent) Description() string {
	return "Returns the event at an index (defaulting to the current controller event index"
}

func (c *GetEvent) Usage() string {
	return "<serial number> [index]"
}

func (c *GetEvent) Help() {
	fmt.Println("Usage: uhppote-cli [options] get-events <serial number> [index]")
	fmt.Println()
	fmt.Println(" Retrieves the event stored at the supplied index. If a specific index is not provided, the command")
	fmt.Println(" defaults to the the current controller event-index and 'bumps' the controller index to the next.")
	fmt.Println(" event. Fails with an error if the supplied index or current controller index is out of range to facilitate")
	fmt.Println(" use with scripts that 'scrape' the event log.")
	fmt.Println()
	fmt.Println(" NOTE: the event index 'wraps around' at 100000 and should not be used as a primary key")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println("  index          (optional) Event index. Defaults to the current controller event index if not provided.")
	fmt.Println("                            Valid index values are:")
	fmt.Println("                            - a number in the range 'first' to 'last' stored event. The event index 'wraps")
	fmt.Println("                              around' at 100000 so if the 'last' index returned by get-events is less than")
	fmt.Println("                              the 'first' index the event log has wrapped and any index in the range 1..100000")
	fmt.Println("                              is valid")
	fmt.Println("                            - 'first' - retrieves the event corresponding to the event at the 'first' event index")
	fmt.Println("                              returned by get-events")
	fmt.Println("                            - 'last' - retrieves the event corresponding to the event at the 'last' event index")
	fmt.Println("                              returned by get-events")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli get-event 12345678 17")
	fmt.Println()
	fmt.Println("    > 405419896  17   2019-07-24 20:12:43 3922570474   1 true  0")
	fmt.Println()
	fmt.Println("      ~ 405419896           controller serial number")
	fmt.Println("      ~ 17                  event index")
	fmt.Println("      ~ 2019-07-24 20:12:43 event timestamp")
	fmt.Println("      ~ 3922570474          card number or user ID")
	fmt.Println("      ~ 1                   door")
	fmt.Println("      ~ true                access granted")
	fmt.Println("      ~ 0                   swipe result")
	fmt.Println()
	fmt.Println("    uhppote-cli get-event 12345678")
	fmt.Println()
	fmt.Println("    > 405419896  23   2019-07-24 20:31:18 3922570474   1 true  0")
	fmt.Println()
	fmt.Println("    uhppote-cli get-event 12345678 first")
	fmt.Println()
	fmt.Println("    > 405419896  1    2019-07-09 21:00:55 3922570474   1 true  0")
	fmt.Println()
	fmt.Println("    uhppote-cli get-event 12345678 last")
	fmt.Println()
	fmt.Println("    > 405419896  69   2019-08-10 10:28:32 3922570474   1 true  44")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *GetEvent) RequiresConfig() bool {
	return false
}

func (c *GetEvent) getIndex(ctx Context, serialNumber uint32) (uint32, error) {
	if len(flag.Args()) > 2 {
		if clean(flag.Args()[2]) == "first" {
			first, err := ctx.uhppote.GetEvent(serialNumber, 0)
			if err != nil {
				return 0, err
			} else if first == nil {
				return 0, fmt.Errorf("Failed to retrieve 'first' event index")
			}

			return first.Index, nil
		}

		if clean(flag.Args()[2]) == "last" {
			last, err := ctx.uhppote.GetEvent(serialNumber, 0xffffffff)
			if err != nil {
				return 0, err
			} else if last == nil {
				return 0, fmt.Errorf("Failed to retrieve 'last' event index")
			}

			return last.Index, nil
		}

		return getUint32(2, "Missing event index", "Invalid event index: %v")
	}

	index, err := ctx.uhppote.GetEventIndex(serialNumber)
	if err != nil {
		return 1, err
	} else if index == nil {
		return 1, fmt.Errorf("Failed to retrieve controller event index")
	}

	return index.Index + 1, nil
}
