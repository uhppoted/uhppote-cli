package commands

import (
	"fmt"
)

var OpenDoorCmd = OpenDoor{}

type OpenDoor struct {
}

func (c *OpenDoor) Execute(ctx Context) error {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return err
	}

	door, err := getUint32(2, "Missing door ID", "Invalid door ID: %v")
	if err != nil {
		return err
	}

	if door != 1 && door != 2 && door != 3 && door != 4 {
		return fmt.Errorf("Invalid door ID: %v", door)
	}

	opened, err := ctx.uhppote.OpenDoor(serialNumber, uint8(door))

	if err == nil {
		fmt.Printf("%v\n", opened)
	}

	return err
}

func (c *OpenDoor) CLI() string {
	return "open"
}

func (c *OpenDoor) Description() string {
	return "Opens a door"
}

func (c *OpenDoor) Usage() string {
	return "<serial number> <door>"
}

func (c *OpenDoor) Help() {
	fmt.Println("Usage: uhppote-cli [options] open <serial number> <door>")
	fmt.Println()
	fmt.Println(" Opens the requested door:")
	fmt.Println()
	fmt.Println("  <serial number>  (required) controller serial number")
	fmt.Println("  <door>           (required) door to open [1,2,3 or 4]")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli open 12345678 2")
	fmt.Println()
}
