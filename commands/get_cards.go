package commands

import (
	"fmt"
)

var GetCardsCmd = GetCards{}

type GetCards struct {
}

func (c *GetCards) Execute(ctx Context) error {
	serialNumber, err := getSerialNumber(ctx)
	if err != nil {
		return err
	}

	N, err := ctx.uhppote.GetCards(serialNumber)
	if err != nil {
		return err
	}

	var index uint32 = 1
	for count := uint32(0); count < N; {
		record, err := ctx.uhppote.GetCardByIndex(serialNumber, index, nil)
		if err != nil {
			return err
		}

		if record != nil {
			fmt.Printf("%v\n", record)
			count++
		}

		index++
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
