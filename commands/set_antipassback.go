package commands

import (
	"flag"
	"fmt"
	"regexp"

	"github.com/uhppoted/uhppote-core/types"
)

var SetAntiPassbackCmd = SetAntiPassback{}

type SetAntiPassback struct {
}

func (c *SetAntiPassback) Execute(ctx Context) error {
	if serialNumber, err := getSerialNumber(ctx); err != nil {
		return err
	} else if antipassback, err := c.parse(); err != nil {
		return err
	} else if ok, err := ctx.uhppote.SetAntiPassback(serialNumber, antipassback); err != nil {
		return err
	} else if ok {
		fmt.Printf("%v  anti-passback %v  ok\n", serialNumber, antipassback)
	} else {
		fmt.Printf("%v  anti-passback %v  failed\n", serialNumber, antipassback)
	}

	return nil
}

func (c *SetAntiPassback) CLI() string {
	return "set-antipassback"
}

func (c *SetAntiPassback) Description() string {
	return "Sets the controller anti-passback mode"
}

func (c *SetAntiPassback) Usage() string {
	return "<serial number> <antipassback>"
}

func (c *SetAntiPassback) Help() {
	fmt.Println("Usage: uhppote-cli [options] set-antipassback <serial-number> <anti-passback> [command options]")
	fmt.Println()
	fmt.Println(" Sets the controller anti-passback mode")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println("  anti-passback  (required) one of the following:")
	fmt.Println("                 - disabled")
	fmt.Println("                 - (1:2);(3:4)")
	fmt.Println("                 - (1:3);(2:4)")
	fmt.Println("                 - 1:(2,3)")
	fmt.Println("                 - 1:(2,3,4)")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    -debug  Displays internal information for diagnosing errors")
	fmt.Println()
	fmt.Println("  Command options:")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *SetAntiPassback) RequiresConfig() bool {
	return false
}

func (c SetAntiPassback) parse() (types.AntiPassback, error) {
	if args := flag.Args(); len(args) > 2 {
		v := regexp.MustCompile(`[ (),]+`).ReplaceAllString(args[2], "")

		switch v {
		case "disabled":
			return types.Disabled, nil

		case "1:2;3:4":
			return types.Readers12_34, nil

		case "1:3;2:4":
			return types.Readers13_24, nil

		case "1:23":
			return types.Readers1_23, nil

		case "1:234":
			return types.Readers1_234, nil

		default:
			return types.Disabled, fmt.Errorf("invalid anti-passback value (%v)", args[2])
		}
	}

	return types.Disabled, fmt.Errorf("missing anti-passback")
}
