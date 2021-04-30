package commands

import (
	"fmt"
	"github.com/uhppoted/uhppoted-api/config"
)

var DeleteAllCmd = DeleteAll{}

type DeleteAll struct {
}

func (c *DeleteAll) Execute(ctx Context) error {
	serialNumber, err := getSerialNumber(ctx)
	if err != nil {
		return err
	}

	deleted, err := ctx.uhppote.DeleteCards(serialNumber)
	if err != nil {
		return err
	}

	fmt.Printf("%v %v\n", serialNumber, deleted)

	return nil
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
	fmt.Printf("              (defaults to %s)\n", config.DefaultConfig)
	fmt.Println("    --debug   Displays internal information for diagnosing errors")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli --debug --config .config delete-all 12345678")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *DeleteAll) RequiresConfig() bool {
	return false
}
