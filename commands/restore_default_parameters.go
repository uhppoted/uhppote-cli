package commands

import (
	"fmt"
	"github.com/uhppoted/uhppoted-lib/config"
)

// RestoreDefaultParameters command variable for CLI command list.
var RestoreDefaultParametersCmd = RestoreDefaultParameters{}

// Command implementation for restore-default-parameters to reset a controller to the manufacturer
// default configuration.
type RestoreDefaultParameters struct {
}

// Gets the controller ID from the command line and sends a restore-default-parameters command to
// the designated controller.
func (c *RestoreDefaultParameters) Execute(ctx Context) error {
	controller, err := getSerialNumber(ctx)
	if err != nil {
		return err
	}

	if ok, err := ctx.uhppote.RestoreDefaultParameters(controller); err != nil {
		return err
	} else if !ok {
		fmt.Printf("%v  restore default parameters failed\n", controller)
	} else {
		fmt.Printf("%v  restore default parameters ok\n", controller)
	}

	return nil
}

// Returns the 'restore-default-parameters' command string for the CLI interface.
func (c *RestoreDefaultParameters) CLI() string {
	return "restore-default-parameters"
}

// Returns the 'restore-default-parameters' command summary for the CLI interface.
func (c *RestoreDefaultParameters) Description() string {
	return "Resets the controller to the manufacturer default configuration"
}

// Returns the 'restore-default-parameters' command parameters for the CLI interface.
func (c *RestoreDefaultParameters) Usage() string {
	return "<serial number>"
}

// Outputs the 'restore-default-parameters' command help for the CLI interface.
func (c *RestoreDefaultParameters) Help() {
	fmt.Println("Usage: uhppote-cli [options] restore-default-parameters <serial number>")
	fmt.Println()
	fmt.Println(" Resets a controller to the manufacturer default configuration")
	fmt.Println()
	fmt.Println("  <serial number>  (required) controller serial number")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --config  File path for the 'conf' file containing the controller configuration")
	fmt.Printf("              (defaults to %s)\n", config.DefaultConfig)
	fmt.Println("    --debug   Displays internal information for diagnosing errors")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli --debug --config .config restore-default-parameters 405419896")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *RestoreDefaultParameters) RequiresConfig() bool {
	return false
}
