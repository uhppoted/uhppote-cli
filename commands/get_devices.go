package commands

import (
	"fmt"
)

var GetDevicesCmd = GetDevices{}

type GetDevices struct {
}

func (c *GetDevices) Execute(ctx Context) error {
	devices, err := ctx.uhppote.FindDevices()

	if err == nil {
		for _, device := range devices {
			fmt.Printf("%s\n", device.String())
		}
	}

	return err
}

func (c *GetDevices) CLI() string {
	return "get-devices"
}

func (c *GetDevices) Description() string {
	return "Returns a list of found UHPPOTE controllers on the network"
}

func (c *GetDevices) Usage() string {
	return ""
}

func (c *GetDevices) Help() {
	fmt.Println("Usage: uhppote-cli [options] get-devices [command options]")
	fmt.Println()
	fmt.Println(" Searches the local network for UHPPOTE access control boards reponding to a poll")
	fmt.Println(" on the default UDP port 60000. Returns a list of boards one per line in the format:")
	fmt.Println()
	fmt.Println(" <serial number> <IP address> <subnet mask> <gateway> <MAC address> <hexadecimal version> <firmware date>")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    -debug  Displays vaguely useful internal information")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *GetDevices) RequiresConfig() bool {
	return false
}
