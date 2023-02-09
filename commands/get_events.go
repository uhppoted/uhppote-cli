package commands

import (
	"fmt"
)

var GetEventsCmd = GetEvents{}

type GetEvents struct {
}

func (c *GetEvents) Execute(ctx Context) error {
	deviceID, err := getSerialNumber(ctx)
	if err != nil {
		return err
	}

	first, err := ctx.uhppote.GetEvent(deviceID, 0)
	if err != nil {
		return err
	}

	last, err := ctx.uhppote.GetEvent(deviceID, 0xffffffff)
	if err != nil {
		return err
	}

	current, err := ctx.uhppote.GetEventIndex(deviceID)
	if err != nil {
		return err
	}

	if first == nil && last == nil {
		fmt.Printf("%v  NO EVENTS\n", deviceID)
	} else if first == nil {
		return fmt.Errorf("failed to get 'first' event")
	} else if last == nil {
		return fmt.Errorf("failed to get 'last' event")
	} else {
		fmt.Printf("%v  %v  %v  %v\n", deviceID, first.Index, last.Index, current.Index)
	}

	return nil
}

func (c *GetEvents) CLI() string {
	return "get-events"
}

func (c *GetEvents) Description() string {
	return "Returns the indices of the 'first' and 'last' events stored on the controller"
}

func (c *GetEvents) Usage() string {
	return "<serial number>"
}

func (c *GetEvents) Help() {
	fmt.Println("Usage: uhppote-cli [options] get-events <serial number>")
	fmt.Println()
	fmt.Println(" Retrieves the indices of the first and last' events stored in the controller event buffer")
	fmt.Println(" The controller event buffer is implemented as a ring buffer with capacity for (apparently)")
	fmt.Println(" 100000 events i.e. the index of the 'last' event may be less than the index of the 'first'")
	fmt.Println(" if the event buffer has wrapped around")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli get-events 12345678")
	fmt.Println()
	fmt.Println("    > 12345678  10  71")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *GetEvents) RequiresConfig() bool {
	return false
}
