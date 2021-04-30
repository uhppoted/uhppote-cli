package commands

import (
	"fmt"
)

var GetCardCmd = GetCard{}

type GetCard struct {
}

func (c *GetCard) Execute(ctx Context) error {
	serialNumber, err := getSerialNumber(ctx)
	if err != nil {
		return err
	}

	cardNumber, err := getUint32(2, "Missing card number", "Invalid card number: %v")
	if err != nil {
		return err
	}

	record, err := ctx.uhppote.GetCardByID(serialNumber, cardNumber)
	if err != nil {
		return err
	}

	if record == nil {
		fmt.Printf("%v %v NO RECORD\n", serialNumber, cardNumber)
	} else {
		fmt.Printf("%-10d %v\n", serialNumber, record)
	}

	return nil
}

func (c *GetCard) CLI() string {
	return "get-card"
}

func (c *GetCard) Description() string {
	return "Returns the access granted to a card number"
}

func (c *GetCard) Usage() string {
	return "<serial number> <card number>"
}

func (c *GetCard) Help() {
	fmt.Println("Usage: uhppote-cli [options] get-card <serial number> <card number>")
	fmt.Println()
	fmt.Println(" Retrieves the access granted for the card number from  the controller card list")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println("  card-number    (required) card number")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli get-card 12345678 9876543")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *GetCard) RequiresConfig() bool {
	return false
}
