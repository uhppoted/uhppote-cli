package commands

import (
	"flag"
	"fmt"
	"github.com/uhppoted/uhppoted-lib/acl"
	"github.com/uhppoted/uhppoted-lib/config"
	"strings"
)

var RevokeCmd = Revoke{}

type Revoke struct {
}

func (c *Revoke) Execute(ctx Context) error {
	if ctx.config == nil {
		return fmt.Errorf("revoke requires a valid configuration file")
	}

	cardNumber, err := getUint32(1, "Missing card number", "Invalid card number: %v")
	if err != nil {
		return err
	}

	doors, err := c.getDoors()
	if err != nil {
		return err
	}

	err = acl.Revoke(ctx.uhppote, ctx.devices, cardNumber, doors)
	if err != nil {
		return err
	}

	fmt.Println("   ... ok")

	return nil
}

func (c *Revoke) getDoors() ([]string, error) {
	doors := []string{}

	s := strings.Join(flag.Args()[2:], " ")
	tokens := strings.Split(s, ",")

	for _, t := range tokens {
		if d := strings.ToLower(strings.ReplaceAll(t, " ", "")); d != "" {
			doors = append(doors, strings.TrimSpace(t))
		}
	}

	return doors, nil
}

func (c *Revoke) CLI() string {
	return "revoke"
}

func (c *Revoke) Description() string {
	return "Revokes access to a door (or doors) for a card "
}

func (c *Revoke) Usage() string {
	return "<card number> <doors>"
}

func (c *Revoke) Help() {
	fmt.Println("Usage: uhppote-cli [options] revoke <card number> <doors>")
	fmt.Println()
	fmt.Println(" Revokes access permissions for a card")
	fmt.Println()
	fmt.Println("  <card number>    (required) card number")
	fmt.Println("  <doors>          (required) comma separated list of permitted doors e.g. Front Door, Workshop")
	fmt.Println("                              Doors are case- and space insensitive and correspond to the doors")
	fmt.Println("                              defined in the config file). 'revoked' permissions are REMOVED from")
	fmt.Println("                              the existing permissions for a card. Use 'delete-card' to remove a")
	fmt.Println("                              card from the internal controller card list. The pseudo-door ALL")
	fmt.Println("                              will revoke the card's access to all doors across all configured")
	fmt.Println("                              controllers.")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --config  File path for the 'conf' file containing the controller configuration")
	fmt.Printf("              (defaults to %s)\n", config.DefaultConfig)
	fmt.Println("    --debug   Displays internal information for diagnosing errors")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli -debug --config .config revoke 918273645 Front Door, Workshop")
	fmt.Println()
}

// Returns true - configuration is not optional for this command to return valid information.
func (c *Revoke) RequiresConfig() bool {
	return true
}
