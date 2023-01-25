package commands

import (
	"flag"
	"fmt"
	"regexp"
	"strings"

	"github.com/uhppoted/uhppoted-lib/config"
)

var SetPCControlCmd = SetPCControl{}

// Command implementation for set-pc-control to enable or disable remote host access control.
//
// The access controller expects the host to communicate at least once every 30 seconds
// otherwise it reverts to local control of access using the stored list of cards (the
// communication does not have to a 'set-pc-control' command). If the access controller
// has reverted to local control because no message has been received from the host for
// more than 30 seconds, any subsequent communication from the remote host will re-establish
// remote control again.
type SetPCControl struct {
}

// Gets the device ID and enable/disable value from the command line and sends a
// set-pc-control command to the designated controller.
func (c *SetPCControl) Execute(ctx Context) error {
	serialNumber, err := getSerialNumber(ctx)
	if err != nil {
		return err
	}

	enable := true
	if len(flag.Args()) > 2 {
		v := strings.ToLower(flag.Arg(2))
		if matches, _ := regexp.MatchString("true|false", v); !matches {
			return fmt.Errorf("Invalid command - expected 'true' or 'false', got '%v'", flag.Arg(2))
		}

		if v == "false" {
			enable = false
		}
	}

	succeeded, err := ctx.uhppote.SetPCControl(serialNumber, enable)
	if err != nil {
		return err
	}

	if !succeeded {
		if enable {
			return fmt.Errorf("Failed to enable 'set PC control' on %v", serialNumber)
		} else {
			return fmt.Errorf("Failed to disable 'set PC control' on %v", serialNumber)
		}
	}

	fmt.Printf("%v %v\n", serialNumber, enable)

	return nil
}

// Returns the 'set-pc-control' command string for the CLI interface.
func (c *SetPCControl) CLI() string {
	return "set-pc-control"
}

// Returns the 'set-pc-control' command summary for the CLI interface.
func (c *SetPCControl) Description() string {
	return "Enables or disables remote host control of access"
}

// Returns the 'set-pc-control' command parameters for the CLI interface.
func (c *SetPCControl) Usage() string {
	return "<serial number> <enabled>"
}

// Outputs the 'set-pc-control' command help for the CLI interface.
func (c *SetPCControl) Help() {
	fmt.Println("Usage: uhppote-cli [options] set-pc-control <serial number> <enable>")
	fmt.Println()
	fmt.Println(" Enables or disables remote host access control.")
	fmt.Println()
	fmt.Println(" The access controller expects a message from the host at least once every")
	fmt.Println(" 30 seconds otherwise it reverts back to local control of access using the")
	fmt.Println(" stored list of cards (the message is not required to be a set-pc-control")
	fmt.Println(" command).")
	fmt.Println()
	fmt.Println("  <serial number>  (required) controller serial number")
	fmt.Println("  <enable>         (optional) 'true' or 'false'. Defaults to 'true'")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --config  File path for the 'conf' file containing the controller configuration")
	fmt.Printf("              (defaults to %s)\n", config.DefaultConfig)
	fmt.Println("    --debug   Displays internal information for diagnosing errors")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli --debug --config .config set-pc-control 12345678")
	fmt.Println("    uhppote-cli --debug --config .config set-pc-control 12345678 true")
	fmt.Println("    uhppote-cli --debug --config .config set-pc-control 12345678 false")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *SetPCControl) RequiresConfig() bool {
	return false
}
