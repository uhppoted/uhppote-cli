package commands

import (
	"fmt"
)

var GetStatusCmd = GetStatus{}

type GetStatus struct {
}

func (c *GetStatus) Execute(ctx Context) error {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return err
	}

	status, err := ctx.uhppote.GetStatus(serialNumber)

	if err == nil {
		fmt.Printf("%v\n", status)
	}

	return err
}

func (c *GetStatus) CLI() string {
	return "get-status"
}

func (c *GetStatus) Description() string {
	return "Returns the current status for the selected controller"
}

func (c *GetStatus) Usage() string {
	return "<serial number>"
}

func (c *GetStatus) Help() {
	fmt.Println("Usage: uhppote-cli [options] get-status <serial number>")
	fmt.Println()
	fmt.Println(" Retrieves the controller status")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli get-status 12345678")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *GetStatus) RequiresConfig() bool {
	return false
}
