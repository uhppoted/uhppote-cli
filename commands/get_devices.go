package commands

import (
	"fmt"
	"sync"
)

var GetDevicesCmd = GetDevices{}

type GetDevices struct {
}

func (c *GetDevices) Execute(ctx Context) error {
	wg := sync.WaitGroup{}
	list := sync.Map{}

	for id, _ := range ctx.uhppote.Devices {
		deviceId := id
		wg.Add(1)
		go func() {
			defer wg.Done()
			if device, err := ctx.uhppote.FindDevice(deviceId); err == nil && device != nil {
				list.Store(deviceId, device.String())
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if devices, err := ctx.uhppote.FindDevices(); err == nil && devices != nil {
			for _, d := range devices {
				list.Store(uint32(d.SerialNumber), d.String())
			}
		}
	}()

	wg.Wait()

	list.Range(func(key, value interface{}) bool {
		fmt.Printf("%v\n", value)
		return true
	})

	return nil
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
	fmt.Println("    --debug   Displays internal information for diagnosing errors")
	fmt.Println("    --config  (optional) configuration file for device specific information")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *GetDevices) RequiresConfig() bool {
	return false
}
