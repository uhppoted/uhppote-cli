package commands

import (
	"flag"
	"fmt"
	"regexp"
	"strings"

	"github.com/uhppoted/uhppoted-api/config"
)

var RecordSpecialEventsCmd = RecordSpecialEvents{}

// Command implementation for record-special-events to enable or disable
// events for door open, door closed and door button pressed.
type RecordSpecialEvents struct {
}

// Gets the device ID and enable/disable value from the command line
// and sends a record-special-events to the designated controller.
func (c *RecordSpecialEvents) Execute(ctx Context) error {
	deviceID, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
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

	succeeded, err := ctx.uhppote.RecordSpecialEvents(deviceID, enable)
	if err != nil {
		return err
	}

	if !succeeded {
		if enable {
			return fmt.Errorf("Failed enable 'record special events' on %v", deviceID)
		} else {
			return fmt.Errorf("Failed disable 'record special events' on %v", deviceID)
		}
	}

	fmt.Printf("%v %v\n", deviceID, enable)

	return nil
}

// Returns the 'record-special-events' command string for the CLI interface.
func (c *RecordSpecialEvents) CLI() string {
	return "record-special-events"
}

// Returns the 'record-special-events' command summary for the CLI interface.
func (c *RecordSpecialEvents) Description() string {
	return "Enables or disables door and pushbutton events"
}

// Returns the 'record-special-events' command parameters for the CLI interface.
func (c *RecordSpecialEvents) Usage() string {
	return "<serial number> <enabled>"
}

// Outputs the 'record-special-events' command help for the CLI interface.
func (c *RecordSpecialEvents) Help() {
	fmt.Println("Usage: uhppote-cli [options] record-special-events <serial number> <enable>")
	fmt.Println()
	fmt.Println(" Enables or disables door and pushbutton events")
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
	fmt.Println("    uhppote-cli --debug --config .config record-special-events 12345678 true")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *RecordSpecialEvents) RequiresConfig() bool {
	return false
}
