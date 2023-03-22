package commands

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/uhppoted/uhppote-core/types"
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

	recordset := []types.Card{}
	index := uint32(1)
	for count := uint32(0); count < N; {
		record, err := ctx.uhppote.GetCardByIndex(serialNumber, index)
		if err != nil {
			c.print(recordset, os.Stdout)
			return err
		}

		if record != nil {
			recordset = append(recordset, *record)
			count++
		}

		index++
	}

	c.print(recordset, os.Stdout)

	return nil
}

func (c *GetCards) print(recordset []types.Card, w io.Writer) error {
	from := func(card types.Card) string {
		if card.From.IsZero() {
			return "-"
		} else {
			return fmt.Sprintf("%v", card.From)
		}
	}

	to := func(card types.Card) string {
		if card.To.IsZero() {
			return "-"
		} else {
			return fmt.Sprintf("%v", card.To)
		}
	}

	door := func(p uint8) string {
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

	pin := func(card types.Card) string {
		if card.PIN > 0 && card.PIN < 1000000 {
			return fmt.Sprintf("%v", card.PIN)
		} else {
			return ""
		}
	}

	table := [][]string{}

	for _, card := range recordset {
		table = append(table, []string{
			fmt.Sprintf("%-8v", card.CardNumber),
			fmt.Sprintf("%-10v", from(card)),
			fmt.Sprintf("%-10v", to(card)),
			fmt.Sprintf("%v", door(card.Doors[1])),
			fmt.Sprintf("%v", door(card.Doors[2])),
			fmt.Sprintf("%v", door(card.Doors[3])),
			fmt.Sprintf("%v", door(card.Doors[4])),
			fmt.Sprintf("%v", pin(card)),
		})
	}

	width := []int{0, 0, 0, 0, 0, 0, 0, 0}
	for _, row := range table {
		for ix, field := range row {
			if len(field) > width[ix] {
				width[ix] = len(field)
			}
		}
	}

	format := fmt.Sprintf("%%-%vv %%%vv %%-%vv %%-%vv %%-%vv %%-%vv %%-%vv %%-%vv\n", width[0], width[1], width[2], width[3], width[4], width[5], width[6], width[7])
	for _, row := range table {
		s := fmt.Sprintf(format, row[0], row[1], row[2], row[3], row[4], row[5], row[6], row[7])
		fmt.Fprintf(w, "%v\n", strings.TrimSpace(s))
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
