package commands

import (
	"fmt"
	"github.com/uhppoted/uhppoted-lib/config"
)

var DeleteCardCmd = DeleteCard{}

type DeleteCard struct {
}

func (c *DeleteCard) Execute(ctx Context) error {
	serialNumber, err := getSerialNumber(ctx)
	if err != nil {
		return err
	}

	cardNumber, err := getUint32(2, "Missing card number", "Invalid card number: %v")
	if err != nil {
		return err
	}

	deleted, err := ctx.uhppote.DeleteCard(serialNumber, cardNumber)
	if err != nil {
		return err
	}

	fmt.Printf("%v %v %v\n", serialNumber, cardNumber, deleted)

	return nil
}

func (c *DeleteCard) CLI() string {
	return "delete-card"
}

func (c *DeleteCard) Description() string {
	return "Deletes a card from the controller"
}

func (c *DeleteCard) Usage() string {
	return "<serial number> <card number>"
}

func (c *DeleteCard) Help() {
	fmt.Println("Usage: uhppote-cli [options] delete-card <serial number> <card number>")
	fmt.Println()
	fmt.Println(" Removes a card from the internal controller card list")
	fmt.Println()
	fmt.Println("  <serial number>  (required) controller serial number")
	fmt.Println("  <card number>    (required) card number")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --config  File path for the 'conf' file containing the controller configuration")
	fmt.Printf("              (defaults to %s)\n", config.DefaultConfig)
	fmt.Println("    --debug   Displays internal information for diagnosing errors")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli --debug --config .config delete-card 12345678 918273645")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *DeleteCard) RequiresConfig() bool {
	return false
}
