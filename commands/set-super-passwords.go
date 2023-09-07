package commands

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/uhppoted/uhppoted-lib/config"
)

// SetSuperCommands variable for CLI command list.
var SetSuperPasswordsCmd = SetSuperPasswords{}

// Command implementation for set-super-password to set up to four 'super' passwords for
// a door.
//
// The command will use up to four of the passwords supplied on the command line. Valid passwords
// are PIN codes in the range [1..999999] and invalid passwords will be replaced with a 0 PIN
// ('no password').
type SetSuperPasswords struct {
}

// Gets the device ID, door and passwords list from the command line and sends a set-super-control
// command to the designated controller.
func (c *SetSuperPasswords) Execute(ctx Context) error {
	controller, err := getSerialNumber(ctx)
	if err != nil {
		return err
	}

	door, err := c.getDoor()
	if err != nil {
		return err
	}

	passwords, err := c.getPasswords()
	if err != nil {
		return err
	}

	if ok, err := ctx.uhppote.SetSuperPasswords(controller, door, passwords...); err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("failed to set super passwords for %v, door %v", controller, door)
	} else {
		fmt.Printf("%v %v\n", controller, ok)
		return nil
	}
}

// Returns the 'set-pc-control' command string for the CLI interface.
func (c *SetSuperPasswords) CLI() string {
	return "set-super-passwords"
}

// Returns the 'set-pc-control' command summary for the CLI interface.
func (c *SetSuperPasswords) Description() string {
	return "Sets the 'super' passwords for a door"
}

// Returns the 'set-pc-control' command parameters for the CLI interface.
func (c *SetSuperPasswords) Usage() string {
	return "<serial number> <door> <passwords>"
}

// Outputs the 'set-super-passwords' command help for the CLI interface.
func (c *SetSuperPasswords) Help() {
	fmt.Println("Usage: uhppote-cli [options] set-super-passwords <serial number> <door> <passwords>")
	fmt.Println()
	fmt.Println(" Sets up to four 'super' passwords for a door.")
	fmt.Println()
	fmt.Println(" Valid passwords are PIN codes in the range [1..999999] and the commands uses the first")
	fmt.Println(" four passwords from the list, replacing invalid passwords with '0' (no password).")
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
	fmt.Println("    uhppote-cli --debug --config .config set-super-passwords 12345678 3 12345,999999,54321")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *SetSuperPasswords) RequiresConfig() bool {
	return false
}

// Returns the door ID from command line argument 3, returning an error if
// missing or not a valid door ID (in the range 1..4]).
func (c *SetSuperPasswords) getDoor() (uint8, error) {
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
func (c *SetSuperPasswords) getPasswords() ([]uint32, error) {
	passwords := []uint32{}

	if len(flag.Args()) > 3 {
		arg := strings.Split(flag.Arg(3), ",")
		if len(arg) > 4 {
			arg = arg[:4]
		}

		re := regexp.MustCompile(`^\s*([0-9]+)\s*$`)

		for _, v := range arg {
			if matches := re.FindStringSubmatch(v); len(matches) > 1 {
				if N, err := strconv.ParseUint(matches[1], 10, 32); err == nil {
					passwords = append(passwords, uint32(N))
				}
			}
		}
	}

	return passwords, nil
}
