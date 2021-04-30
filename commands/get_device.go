package commands

import (
	"fmt"
)

var GetDeviceCmd = GetDevice{}

type GetDevice struct {
}

func (c *GetDevice) Execute(ctx Context) error {
	serialNumber, err := getSerialNumber(ctx)
	if err != nil {
		return err
	}

	device, err := ctx.uhppote.GetDevice(serialNumber)
	if err != nil {
		return err
	} else if device == nil {
		return fmt.Errorf("No device found matching serial number '%d'", serialNumber)
	}

	fmt.Printf("%s\n", device.String())

	return nil
}

func (c *GetDevice) CLI() string {
	return "get-device"
}

func (c *GetDevice) Description() string {
	return "'pings' a UHPPOTE controller using the IP address configured for the device"
}

func (c *GetDevice) Usage() string {
	return "<serial number>"
}

func (c *GetDevice) Help() {
	fmt.Println("Usage: uhppote-cli [options] get-device <serial number>>")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println()
	fmt.Println(" Issues a 'get-devices' request directed at the IP address configured for the supplied serial number")
	fmt.Println(" and extracts the response matching the serial number, returning the board summary information in the format:")
	fmt.Println()
	fmt.Println(" <serial number> <IP address> <subnet mask> <gateway> <MAC address> <hexadecimal version> <firmware date>")
	fmt.Println()
	fmt.Println(" Falls back to a broadcast on the local network if no IP address is configured for the supplied serial number.")
	fmt.Println()
	fmt.Println("  Example:")
	fmt.Println()
	fmt.Println("    uhppote-cli get-device 12345678")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    -debug  Displays a trace of request/response messages")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *GetDevice) RequiresConfig() bool {
	return false
}
