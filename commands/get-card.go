package commands

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/uhppoted/uhppote-core/types"
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
		c.println(*record, os.Stdout)
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
	fmt.Println("    uhppote-cli get-card 405419896 10058400")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *GetCard) RequiresConfig() bool {
	return false
}

func (c GetCard) println(card types.Card, w io.Writer) error {
	var b bytes.Buffer

	if err := c.print(card, &b); err != nil {
		return err
	}

	fmt.Fprintf(w, "%v\n", b.String())

	return nil
}

func (c GetCard) print(card types.Card, w io.Writer) error {
	f := func(p uint8) string {
		switch {
		case p == 0:
			return "N"

		case p == 1:
			return "Y"

		case p >= 2 && p <= 254:
			return fmt.Sprintf("%v", p)

		default:
			return "N"
		}
	}

	from := "-"
	to := "-"
	doors := fmt.Sprintf("%v %v %v %v", f(card.Doors[1]), f(card.Doors[2]), f(card.Doors[3]), f(card.Doors[4]))
	PIN := "-"
	firstcard := "-"

	if !card.From.IsZero() {
		from = fmt.Sprintf("%v", card.From)
	}

	if !card.To.IsZero() {
		to = fmt.Sprintf("%v", card.To)
	}

	if card.PIN != 0 && card.PIN <= 999999 {
		PIN = fmt.Sprintf("%v", card.PIN)
	}

	if !card.FirstCard.IsZero() {
		firstcard = fmt.Sprintf("%v", card.FirstCard)
	}

	fmt.Fprintf(w, "%-8v %-10v %-10v %v %v %v", card.CardNumber, from, to, doors, PIN, firstcard)

	return nil
}
