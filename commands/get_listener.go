package commands

import (
	"fmt"
)

var GetListenerCmd = GetListener{}

type GetListener struct {
}

func (c *GetListener) Execute(ctx Context) error {
	serialNumber, err := getSerialNumber(ctx)
	if err != nil {
		return err
	}

	if address, interval, err := ctx.uhppote.GetListener(serialNumber); err != nil {
		return err
	} else if interval > 0 {
		fmt.Printf("%v %v %vs\n", serialNumber, address, interval)
	} else {
		fmt.Printf("%v %v\n", serialNumber, address)
	}

	return nil
}

func (c *GetListener) CLI() string {
	return "get-listener"
}

func (c *GetListener) Description() string {
	return "Returns the IPv4 address:port to which the controller sends events"
}

func (c *GetListener) Usage() string {
	return "<serial number>"
}

func (c *GetListener) Help() {
	fmt.Println("Usage: uhppote-cli [options] get-listener <serial number>")
	fmt.Println()
	fmt.Println(" Retrieves the configured IP address:port to which the controller sends events. Also")
	fmt.Println(" retrieves the controller auto-send interval (if not zero).")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println()
	fmt.Println("  Example:")
	fmt.Println()
	fmt.Println("    uhppote-cli get-listener 405419896")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *GetListener) RequiresConfig() bool {
	return false
}
