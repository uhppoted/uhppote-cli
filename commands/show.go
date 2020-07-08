package commands

import (
	"fmt"
	"github.com/uhppoted/uhppoted-api/acl"
	"github.com/uhppoted/uhppoted-api/config"
	"sort"
	"strings"
)

var ShowCmd = Show{}

type Show struct {
}

func (c *Show) Execute(ctx Context) error {
	if ctx.config == nil {
		return fmt.Errorf("show requires a valid configuration file")
	}

	cardNumber, err := getUint32(1, "Missing card number", "Invalid card number: %v")
	if err != nil {
		return err
	}

	devices := getDevices(&ctx)

	permissions, err := acl.GetCard(ctx.uhppote, devices, cardNumber)
	if err != nil {
		return err
	}

	doors := []string{}
	width := 0
	for k := range permissions {
		doors = append(doors, k)

		if width < len([]rune(k)) {
			width = len([]rune(k))
		}
	}

	sort.Slice(doors, func(i, j int) bool {
		p := strings.ToLower(strings.ReplaceAll(doors[i], " ", ""))
		q := strings.ToLower(strings.ReplaceAll(doors[j], " ", ""))
		return p < q
	})

	fmt.Println()
	format := fmt.Sprintf("%%-%ds  %%v  %%v\n", width)
	for _, door := range doors {
		v, _ := permissions[door]
		fmt.Printf(format, door, v.From, v.To)
	}
	fmt.Println()

	return nil
}

func (c *Show) CLI() string {
	return "show"
}

func (c *Show) Description() string {
	return "Lists the access permissions for a card"
}

func (c *Show) Usage() string {
	return "<card number>"
}

func (c *Show) Help() {
	fmt.Println("Usage: uhppote-cli [options] show <card number>")
	fmt.Println()
	fmt.Println(" Lists the access permissions for a card")
	fmt.Println()
	fmt.Println("  <card number>    (required) card number")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --config  File path for the 'conf' file containing the controller configuration")
	fmt.Printf("              (defaults to %s)\n", config.DefaultConfig)
	fmt.Println("    --debug   Displays internal information for diagnosing errors")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli show 918273645")
	fmt.Println()
}

// Returns true - configuration is not optional for this command to return valid information.
func (c *Show) RequiresConfig() bool {
	return true
}
