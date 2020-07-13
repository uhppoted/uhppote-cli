package commands

import (
	"fmt"
)

var GetEventIndexCmd = GetEventIndex{}

type GetEventIndex struct {
}

func (c *GetEventIndex) Execute(ctx Context) error {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return err
	}

	index, err := ctx.uhppote.GetEventIndex(serialNumber)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", index)

	return nil
}

func (c *GetEventIndex) CLI() string {
	return "get-event-index"
}

func (c *GetEventIndex) Description() string {
	return "Retrieves the current event index"
}

func (c *GetEventIndex) Usage() string {
	return "<serial number>"
}

func (c *GetEventIndex) Help() {
	fmt.Println("Usage: uhppote-cli [options] get-event-index <serial number>")
	fmt.Println()
	fmt.Println(" Retrieves the current event record index")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli get-event-index 12345678")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *GetEventIndex) RequiresConfig() bool {
	return false
}
