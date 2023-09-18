package commands

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/uhppoted/uhppoted-lib/config"
)

// SetDoorPasscodes command variable for CLI command list.
var SetDoorPasscodesCmd = SetDoorPasscodes{}

// Command implementation for set-door-passcodes to set up to four supervisor passcodes for
// a door.
//
// The command will use up to four of the codes supplied on the command line. Valid passcodes
// are PIN codes in the range [1..999999] and invalid codes will be replaced with a 0 PIN
// ('no code').
type SetDoorPasscodes struct {
}

// Gets the device ID, door and passwords list from the command line and sends a set-super-control
// command to the designated controller.
func (c *SetDoorPasscodes) Execute(ctx Context) error {
	controller, err := getSerialNumber(ctx)
	if err != nil {
		return err
	}

	door, err := c.getDoor()
	if err != nil {
		return err
	}

	passcodes, err := c.getPasscodes()
	if err != nil {
		return err
	}

	if ok, err := ctx.uhppote.SetDoorPasscodes(controller, door, passcodes...); err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("failed to set door passcodes for %v, door %v", controller, door)
	} else {
		fmt.Printf("%v %v\n", controller, ok)
		return nil
	}
}

// Returns the 'set-door-passcodes' command string for the CLI interface.
func (c *SetDoorPasscodes) CLI() string {
	return "set-door-passcodes"
}

// Returns the 'set-door-passcodes' command summary for the CLI interface.
func (c *SetDoorPasscodes) Description() string {
	return "Sets the supervisor passcodes for a door"
}

// Returns the 'set-door-passcodes' command parameters for the CLI interface.
func (c *SetDoorPasscodes) Usage() string {
	return "<serial number> <door> <passcodes>"
}

// Outputs the 'set-door-passcodes' command help for the CLI interface.
func (c *SetDoorPasscodes) Help() {
	fmt.Println("Usage: uhppote-cli [options] set-door-passcodes <serial number> <door> <passcodes>")
	fmt.Println()
	fmt.Println(" Sets up to four supervisor passcodes for a door.")
	fmt.Println()
	fmt.Println(" Valid passcodes are PIN codes in the range [1..999999] and the commands uses the first")
	fmt.Println(" four codes from the list, replacing invalid passcodes with '0' (no code).")
	fmt.Println()
	fmt.Println("  <serial number>  (required) controller serial number")
	fmt.Println("  <door>           (required) door [1..4]")
	fmt.Println("  <passwords>      (required) comma seperated list of passwords")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --config  File path for the 'conf' file containing the controller configuration")
	fmt.Printf("              (defaults to %s)\n", config.DefaultConfig)
	fmt.Println("    --debug   Displays internal information for diagnosing errors")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli --debug --config .config set-door-passcodes 12345678 3 12345,999999,54321")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *SetDoorPasscodes) RequiresConfig() bool {
	return false
}

// Returns the door ID from command line argument 3, returning an error if
// missing or not a valid door ID (in the range 1..4]).
func (c *SetDoorPasscodes) getDoor() (uint8, error) {
	if len(flag.Args()) < 3 {
		return 0, fmt.Errorf("missing door ID")
	}

	arg := flag.Arg(2)

	if valid, _ := regexp.MatchString("[1-4]", arg); !valid {
		return 0, fmt.Errorf("invalid door ID (%v)", arg)
	}

	if N, err := strconv.ParseUint(arg, 10, 8); err != nil {
		return 0, fmt.Errorf("invalid door ID (%v)", arg)
	} else {
		return uint8(N), nil
	}
}

// Returns a list of up to four passwords converted from command line argument 4.
func (c *SetDoorPasscodes) getPasscodes() ([]uint32, error) {
	passcodes := []uint32{}

	if len(flag.Args()) > 3 {
		arg := strings.Split(flag.Arg(3), ",")
		if len(arg) > 4 {
			arg = arg[:4]
		}

		re := regexp.MustCompile(`^\s*([0-9]+)\s*$`)

		for _, v := range arg {
			if matches := re.FindStringSubmatch(v); len(matches) > 1 {
				if N, err := strconv.ParseUint(matches[1], 10, 32); err == nil {
					passcodes = append(passcodes, uint32(N))
				}
			}
		}
	}

	return passcodes, nil
}
