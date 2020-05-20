package commands

import (
	"fmt"
)

var GetTimeCmd = GetTime{}

type GetTime struct {
}

func (c *GetTime) Execute(ctx Context) error {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return err
	}

	datetime, err := ctx.uhppote.GetTime(serialNumber)

	if err == nil {
		fmt.Printf("%v\n", datetime)
	}

	return err
}

func (c *GetTime) CLI() string {
	return "get-time"
}

func (c *GetTime) Description() string {
	return "Returns the current time on the selected controller"
}

func (c *GetTime) Usage() string {
	return "<serial number>"
}

func (c *GetTime) Help() {
	fmt.Println("Usage: uhppote-cli [options] get-time <serial number> [command options]")
	fmt.Println()
	fmt.Println(" Retrieves the current date/time referenced to the local timezone for the controller")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    -debug  Displays vaguely useful internal information")
	fmt.Println()
	fmt.Println("  Command options:")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *GetTime) RequiresConfig() bool {
	return false
}
