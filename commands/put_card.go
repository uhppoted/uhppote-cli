package commands

import (
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-api/config"
)

var PutCardCmd = PutCard{}

type PutCard struct {
}

func (c *PutCard) Execute(ctx Context) error {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return err
	}

	cardNumber, err := getUint32(2, "Missing card number", "Invalid card number: %v")
	if err != nil {
		return err
	}

	from, err := getDate(3, "Missing start date", "Invalid start date: %v")
	if err != nil {
		return err
	}

	to, err := getDate(4, "Missing end date", "Invalid end date: %v")
	if err != nil {
		return err
	}

	permissions, err := getPermissions(5)
	if err != nil {
		return err
	}

	start := types.Date(*from)
	end := types.Date(*to)

	authorised, err := ctx.uhppote.PutCard(serialNumber, types.Card{
		CardNumber: cardNumber,
		From:       &start,
		To:         &end,
		Doors:      permissions,
	})

	if err != nil {
		return err
	}

	fmt.Printf("%v\n", authorised)
	return nil
}

func (c *PutCard) CLI() string {
	return "put-card"
}

func (c *PutCard) Description() string {
	return "Sets the card access rights on a single access controller"
}

func (c *PutCard) Usage() string {
	return "<serial number> <card number> <start date> <end date> <doors>"
}

func (c *PutCard) Help() {
	fmt.Println("Usage: uhppote-cli [options] put-card <serial number> <card number> <start date> <end date> <doors>")
	fmt.Println()
	fmt.Println(" Adds (or updates) a card to the list of the cards managed by a controller")
	fmt.Println()
	fmt.Println("  <serial number>  (required) controller serial number")
	fmt.Println("  <card number>    (required) card number")
	fmt.Println("  <start date>     (required) start date YYYY-MM-DD")
	fmt.Println("  <end date>       (required) end date   YYYY-MM-DD")
	fmt.Println("  <doors>          (required) list of permitted doors [1,2,3,4]. Unlisted doors will be set to 'N'")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --config  File path for the 'conf' file containing the controller configuration")
	fmt.Printf("              (defaults to %s)\n", config.DefaultConfig)
	fmt.Println("    --debug   Displays internal information for diagnosing errors")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli put-card 12345678 918273645 2020-01-01 2020-12-31 1,2,4")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *PutCard) RequiresConfig() bool {
	return false
}
