package commands

import (
	"fmt"
)

var DeleteCardCmd = DeleteCard{}

type DeleteCard struct {
}

func (c *DeleteCard) Execute(ctx Context) error {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return err
	}

	cardNumber, err := getUint32(2, "Missing card number", "Invalid card number: %v")
	if err != nil {
		return err
	}

	result, err := ctx.uhppote.DeleteCard(serialNumber, cardNumber)

	if err == nil {
		fmt.Printf("%v\n", result)
	}

	return err
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
	fmt.Printf("              (defaults to %s)\n", DefaultConfig)
	fmt.Println("    --debug   Displays vaguely useful internal information")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli --debug --config .config delete-card 12345678 918273645")
	fmt.Println()
}
