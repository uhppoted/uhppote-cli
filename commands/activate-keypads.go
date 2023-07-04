package commands

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

var ActivateKeypadsCmd = ActivateKeypads{}

type ActivateKeypads struct {
}

func (c *ActivateKeypads) Execute(ctx Context) error {
	serialNumber, err := getSerialNumber(ctx)
	if err != nil {
		return err
	}

	keypads, err := c.getKeypads()
	if err != nil {
		return err
	}

	if activated, err := ctx.uhppote.ActivateKeypads(serialNumber, keypads); err != nil {
		return err
	} else if !activated {
		return fmt.Errorf("failed to activate access keypads")
	} else {
		readers := []string{}

		for _, reader := range []uint8{1, 2, 3, 4} {
			if keypads[reader] {
				readers = append(readers, fmt.Sprintf("%v", reader))
			}
		}

		if len(readers) == 0 {
			fmt.Printf("%v  activated keypads %v\n", serialNumber, "(none)")
		} else {
			fmt.Printf("%v  activated keypads %v\n", serialNumber, strings.Join(readers, ","))
		}
	}

	return nil
}

func (c *ActivateKeypads) CLI() string {
	return "activate-keypads"
}

func (c *ActivateKeypads) Description() string {
	return "Activates (or deactivates) the access keypads on a controller"
}

func (c *ActivateKeypads) Usage() string {
	return "<serial number> <keypads>"
}

func (c *ActivateKeypads) Help() {
	fmt.Println("Usage: uhppote-cli [options] activate-keypads <serial number> <doors>")
	fmt.Println()
	fmt.Println(" Activates the keypads for the doors (unlisted keypads are deactivated)")
	fmt.Println()
	fmt.Println("  <serial number>  (required) controller serial number")
	fmt.Println("  <doors>          (required) comma seperated list of doors with keypads e.g. 1,2,4")
	fmt.Println()
	fmt.Println("  Example:")
	fmt.Println()
	fmt.Println("    uhppote-cli activate-keypads 405419896 1,4")
	fmt.Println()
	fmt.Println("    (activates keypads on door 1 and 4, deactivates keypads on doors 2 and 3)")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *ActivateKeypads) RequiresConfig() bool {
	return false
}

func (c *ActivateKeypads) getKeypads() (map[uint8]bool, error) {
	keypads := map[uint8]bool{
		1: false,
		2: false,
		3: false,
		4: false,
	}

	if args := flag.Args(); len(args) > 2 {
		doors := strings.Split(args[2], ",")
		for _, d := range doors {
			if door, err := strconv.Atoi(d); err != nil {
				return nil, err
			} else if door < 1 || door > 4 {
				return nil, fmt.Errorf("invalid door (%v)", d)
			} else {
				keypads[uint8(door)] = true
			}
		}
	}

	return keypads, nil
}
