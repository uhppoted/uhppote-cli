package commands

import (
	"fmt"
)

var SetEventIndexCmd = SetEventIndex{}

type SetEventIndex struct {
}

func (c *SetEventIndex) Execute(ctx Context) error {
	serialNumber, err := getSerialNumber(ctx)
	if err != nil {
		return err
	}

	index, err := getUint32(2, "Missing event index", "Invalid event index: %v")
	if err != nil {
		return err
	}

	result, err := ctx.uhppote.SetEventIndex(serialNumber, index)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", result)

	return nil
}

func (c *SetEventIndex) CLI() string {
	return "set-event-index"
}

func (c *SetEventIndex) Description() string {
	return "Sets the current event index"
}

func (c *SetEventIndex) Usage() string {
	return "<serial number> <index>"
}

func (c *SetEventIndex) Help() {
	fmt.Println("Usage: uhppote-cli [options] set-event-index <serial number> <index>")
	fmt.Println()
	fmt.Println(" Sets the event index")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println("  index          (required) event index")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli set-event-index 12345678 15")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *SetEventIndex) RequiresConfig() bool {
	return false
}
