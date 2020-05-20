package commands

import (
	"fmt"
)

var GetCardsCmd = GetCards{}

type GetCards struct {
}

func (c *GetCards) Execute(ctx Context) error {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return err
	}

	N, err := ctx.uhppote.GetCards(serialNumber)

	if err != nil {
		return err
	}

	for index := uint32(0); index < N.Records; index++ {
		record, err := ctx.uhppote.GetCardByIndex(serialNumber, index+1)
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", record)
	}

	return nil
}

func (c *GetCards) CLI() string {
	return "get-cards"
}

func (c *GetCards) Description() string {
	return "Returns the list of cards stored on the controller"
}

func (c *GetCards) Usage() string {
	return "<serial number>"
}

func (c *GetCards) Help() {
	fmt.Println("Usage: uhppote-cli [options] get-cards <serial number>")
	fmt.Println()
	fmt.Println(" Retrieves the number of cards in the controller card list")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli get-cards 12345678")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *GetCards) RequiresConfig() bool {
	return false
}
