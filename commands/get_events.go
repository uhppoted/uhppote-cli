package commands

import (
	"errors"
	"fmt"
)

var GetEventsCmd = GetEvents{}

type GetEvents struct {
}

func (c *GetEvents) Execute(ctx Context) error {
	serialNumber, err := getSerialNumber(ctx)
	if err != nil {
		return err
	}

	first, err := ctx.uhppote.GetEvent(serialNumber, 0)
	if err != nil {
		return err
	}

	last, err := ctx.uhppote.GetEvent(serialNumber, 0xffffffff)
	if err != nil {
		return err
	}

	if first == nil && last == nil {
		fmt.Printf("%v  NO EVENTS\n", serialNumber)
	} else if first == nil {
		return errors.New("Failed to get 'first' event")
	} else if last == nil {
		return errors.New("Failed to get 'last' event")
	} else {
		fmt.Printf("%v  %d  %d\n", serialNumber, first.Index, last.Index)
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
