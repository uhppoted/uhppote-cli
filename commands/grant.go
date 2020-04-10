package commands

import (
	"flag"
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-api/acl"
	"strings"
)

var GRANT = Grant{}

type Grant struct {
}

func (c *Grant) Execute(ctx Context) error {
	if ctx.config == nil {
		return fmt.Errorf("grant requires a valid configuration file")
	}

	err := ctx.config.Verify()
	if err != nil {
		return err
	}

	cardNumber, err := getUint32(1, "Missing card number", "Invalid card number: %v")
	if err != nil {
		return err
	}

	from, err := getDate(2, "Missing start date", "Invalid start date: %v")
	if err != nil {
		return err
	}

	to, err := getDate(3, "Missing end date", "Invalid end date: %v")
	if err != nil {
		return err
	}

	doors, err := c.getDoors()
	if err != nil {
		return err
	}

	devices := getDevices(&ctx)

	return acl.Grant(ctx.uhppote, devices, cardNumber, types.Date(*from), types.Date(*to), doors)
}

func (c *Grant) getDoors() ([]string, error) {
	doors := []string{}

	s := strings.Join(flag.Args()[4:], " ")
	tokens := strings.Split(s, ",")

	for _, t := range tokens {
		if d := strings.ToLower(strings.ReplaceAll(t, " ", "")); d != "" {
			doors = append(doors, d)
		}
	}

	return doors, nil
}

func (c *Grant) CLI() string {
	return "grant"
}

func (c *Grant) Description() string {
	return "Grants a card access to a door (or doors)"
}

func (c *Grant) Usage() string {
	return "<card number> <start date> <end date> <doors>"
}

func (c *Grant) Help() {
	fmt.Println("Usage: uhppote-cli [options] grant <card number> <start date> <end date> <doors>")
	fmt.Println()
	fmt.Println(" Sets the access permissions for a card")
	fmt.Println()
	fmt.Println("  <card number>    (required) card number")
	fmt.Println("  <start date>     (required) start date YYYY-MM-DD")
	fmt.Println("  <end date>       (required) end date   YYYY-MM-DD")
	fmt.Println("  <doors>          (required) comma separated list of permitted doors e.g. Front Door, Workshop")
	fmt.Println("                              Doors are case- and space insensitive and correspond to the doors")
	fmt.Println("                              defined in the config file.")
	fmt.Println()
	fmt.Println("                              N.B. 'grant' permissions are ADDED to the existing permissions for")
	fmt.Println("                                    a card. Use 'revoke' to remove unwanted permissions.")
	fmt.Println("                                    Also, the 'from' and 'to' dates for a card are WIDENED to")
	fmt.Println("                                    the earliest 'from' date and latest 'to' date combination")
	fmt.Println("                                    for all records for this card across all controllers.")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --config  File path for the 'conf' file containing the controller configuration")
	fmt.Printf("              (defaults to %s)\n", DEFAULT_CONFIG)
	fmt.Println("    --debug   Displays vaguely useful internal information")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli grant 918273645 2020-01-01 2020-12-31 Front Door, Workshop")
	fmt.Println()
}
