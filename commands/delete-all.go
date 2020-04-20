package commands

import (
	"fmt"
)

var DeleteAllCmd = DeleteAll{}

type DeleteAll struct {
}

func (c *DeleteAll) Execute(ctx Context) error {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return err
	}

	deleted, err := ctx.uhppote.DeleteCards(serialNumber)

	if err == nil {
		fmt.Printf("%v\n", deleted)
	}

	return err
}

func (c *DeleteAll) CLI() string {
	return "delete-all"
}

func (c *DeleteAll) Description() string {
	return "Clears all cards stored on the controller"
}

func (c *DeleteAll) Usage() string {
	return "<serial number>"
}

func (c *DeleteAll) Help() {
	fmt.Println("Usage: uhppote-cli [options] delete-all <serial number>")
	fmt.Println()
	fmt.Println(" Removes all cards from the controller internal card list")
	fmt.Println()
	fmt.Println("  <serial number>  (required) controller serial number")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --config  File path for the 'conf' file containing the controller configuration")
	fmt.Printf("              (defaults to %s)\n", DefaultConfig)
	fmt.Println("    --debug   Displays vaguely useful internal information")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli --debug --config .config delete-all 12345678")
	fmt.Println()
}
