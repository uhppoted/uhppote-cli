package commands

import (
	"fmt"

	"github.com/uhppoted/uhppoted-lib/config"
)

var GetTimeProfileCmd = GetTimeProfile{}

type GetTimeProfile struct {
}

func (c *GetTimeProfile) Execute(ctx Context) error {
	serialNumber, err := getSerialNumber(ctx)
	if err != nil {
		return err
	}

	profileID, err := getUint8(2, "missing time profile ID", "invalid time profile ID: %v")
	if err != nil {
		return err
	} else if profileID < 2 || profileID > 254 {
		return fmt.Errorf("invalid time profile ID (%v) - valid range is from 2 to 254", profileID)
	}

	profile, err := ctx.uhppote.GetTimeProfile(serialNumber, profileID)
	if err != nil {
		return err
	}

	if profile == nil {
		fmt.Printf("%v %v NO ACTIVE TIME PROFILE\n", serialNumber, profileID)
	} else {
		fmt.Printf("%-10d %v\n", serialNumber, profile)
	}

	return nil
}

func (c *GetTimeProfile) CLI() string {
	return "get-time-profile"
}

func (c *GetTimeProfile) Description() string {
	return "Retrieves the time profile associated with a time profile ID"
}

func (c *GetTimeProfile) Usage() string {
	return "<serial number> <profile ID>"
}

func (c *GetTimeProfile) Help() {
	fmt.Println("Usage: uhppote-cli [options] get-time-profile <serial number> <profile ID>")
	fmt.Println()
	fmt.Println(" Retrieves the time profile associated with a profile ID")
	fmt.Println()
	fmt.Println("  serial number  (required) controller serial number")
	fmt.Println("  profile ID     (required) time profile ID (2-254)")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --config  File path for the 'conf' file containing the controller configuration")
	fmt.Printf("              (defaults to %s)\n", config.DefaultConfig)
	fmt.Println("    --debug   Displays internal information for diagnosing errors")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli get-time-profile 9876543210 7")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *GetTimeProfile) RequiresConfig() bool {
	return false
}
