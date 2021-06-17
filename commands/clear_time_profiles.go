package commands

import (
	"fmt"

	"github.com/uhppoted/uhppoted-lib/config"
)

var ClearTimeProfilesCmd = ClearTimeProfiles{}

type ClearTimeProfiles struct {
}

func (c *ClearTimeProfiles) Execute(ctx Context) error {
	serialNumber, err := getSerialNumber(ctx)
	if err != nil {
		return err
	}

	cleared, err := ctx.uhppote.ClearTimeProfiles(serialNumber)
	if err != nil {
		return err
	}

	fmt.Printf("%v %v\n", serialNumber, cleared)

	return nil
}

func (c *ClearTimeProfiles) CLI() string {
	return "clear-time-profiles"
}

func (c *ClearTimeProfiles) Description() string {
	return "Clears all time profiles on the controller"
}

func (c *ClearTimeProfiles) Usage() string {
	return "<serial number>"
}

func (c *ClearTimeProfiles) Help() {
	fmt.Println("Usage: uhppote-cli [options] clear-time-profiles <serial number>")
	fmt.Println()
	fmt.Println(" Clears all time profiles from the controller")
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
	fmt.Println("    uhppote-cli --debug clear-time-profiles 9876543210")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *ClearTimeProfiles) RequiresConfig() bool {
	return false
}
