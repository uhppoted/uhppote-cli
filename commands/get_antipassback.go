package commands

import (
	"fmt"
)

var GetAntiPassbackCmd = GetAntiPassback{}

type GetAntiPassback struct {
}

func (c *GetAntiPassback) Execute(ctx Context) error {
	if serialNumber, err := getSerialNumber(ctx); err != nil {
		return err
	} else if antipassback, err := ctx.uhppote.GetAntiPassback(serialNumber); err != nil {
		return err
	} else {
		fmt.Printf("anti-passback: %v\n", antipassback)
	}

	return nil
}

func (c *GetAntiPassback) CLI() string {
	return "get-antipassback"
}

func (c *GetAntiPassback) Description() string {
	return "Returns the controller anti-passback setting"
}

func (c *GetAntiPassback) Usage() string {
	return "<serial number>"
}

func (c *GetAntiPassback) Help() {
	fmt.Println("Usage: uhppote-cli [options] get-antipassback <serial number> [command options]")
	fmt.Println()
	fmt.Println(" Retrieves the controller anti-passback setting")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    -debug  Displays internal information for diagnosing errors")
	fmt.Println()
	fmt.Println("  Command options:")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *GetAntiPassback) RequiresConfig() bool {
	return false
}
