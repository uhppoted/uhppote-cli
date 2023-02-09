package commands

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-lib/acl"
	"github.com/uhppoted/uhppoted-lib/config"
)

var GrantCmd = Grant{}

type Grant struct {
}

func (c *Grant) Execute(ctx Context) error {
	if ctx.config == nil {
		return fmt.Errorf("grant requires a valid configuration file")
	}

	cardNumber, err := getUint32(1, "missing card number", "invalid card number: %v")
	if err != nil {
		return err
	}

	from, err := getDate(2, "missing start date", "invalid start date: %v")
	if err != nil {
		return err
	}

	to, err := getDate(3, "missing end date", "invalid end date: %v")
	if err != nil {
		return err
	}

	var re = regexp.MustCompile("[0-9]+")
	var profileID = 0
	var doors []string

	if len(flag.Args()) > 5 && re.MatchString(flag.Arg(4)) {
		profileID, err = strconv.Atoi(flag.Arg(4))
		if err != nil {
			return err
		} else if profileID < 2 || profileID > 254 {
			return fmt.Errorf("invalid time profile ID (%v) - valid range is from 2 to 254", profileID)
		}

		doors, err = c.getDoors(5)
		if err != nil {
			return err
		}
	} else if doors, err = c.getDoors(4); err != nil {
		return err
	}

	err = acl.Grant(ctx.uhppote, ctx.devices, cardNumber, types.Date(*from), types.Date(*to), profileID, doors)
	if err != nil {
		return err
	}

	fmt.Println(" ... ok")

	return nil
}

func (c *Grant) getDoors(ix int) ([]string, error) {
	doors := []string{}

	s := strings.Join(flag.Args()[ix:], " ")
	tokens := strings.Split(s, ",")

	for _, t := range tokens {
		if d := strings.ToLower(strings.ReplaceAll(t, " ", "")); d != "" {
			doors = append(doors, strings.TrimSpace(t))
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
	return "<card number> <start date> <end date> [profile] <doors>"
}

func (c *Grant) Help() {
	fmt.Println("Usage: uhppote-cli [options] grant <card number> <start date> <end date> <profile> <doors>")
	fmt.Println()
	fmt.Println(" Sets the access permissions for a card")
	fmt.Println()
	fmt.Println("  <card number>    (required) card number")
	fmt.Println("  <start date>     (required) start date YYYY-MM-DD")
	fmt.Println("  <end date>       (required) end date   YYYY-MM-DD")
	fmt.Println("  <profile>        (optional) predefined time profile, in the range [2..254]")
	fmt.Println("  <doors>          (required) comma separated list of permitted doors e.g. Front Door, Workshop")
	fmt.Println("                              Doors are case- and space insensitive and correspond to the doors")
	fmt.Println("                              defined in the config file. The pseudo-door ALL will grant the")
	fmt.Println("                              card access to all doors across all configured devices")
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
	fmt.Printf("              (defaults to %s)\n", config.DefaultConfig)
	fmt.Println("    --debug   Displays internal information for diagnosing errors")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli grant 918273645 2020-01-01 2020-12-31 Front Door, Workshop")
	fmt.Println(`    uhppote-cli grant 918273645 2020-01-01 2020-12-31 29 "Front Door, Workshop"`)
	fmt.Println("    uhppote-cli grant 918273645 2020-01-01 2020-12-31 ALL")
	fmt.Println()
}

// Returns true - configuration is not optional for this command to return valid information.
func (c *Grant) RequiresConfig() bool {
	return true
}
