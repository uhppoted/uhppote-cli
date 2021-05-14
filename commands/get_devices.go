package commands

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/uhppoted/uhppote-core/types"
)

var GetDevicesCmd = GetDevices{}

type GetDevices struct {
}

func (c *GetDevices) Execute(ctx Context) error {
	wg := sync.WaitGroup{}
	list := sync.Map{}

	for id, _ := range ctx.uhppote.DeviceList() {
		deviceId := id
		wg.Add(1)
		go func() {
			defer wg.Done()
			if device, err := ctx.uhppote.GetDevice(deviceId); err != nil {
				fmt.Fprintf(os.Stderr, "   WARN:  %v\n", err)
			} else if device != nil {
				list.Store(deviceId, *device)
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if devices, err := ctx.uhppote.GetDevices(); err != nil {
			fmt.Fprintf(os.Stderr, "   WARN:  %v\n", err)
		} else if devices != nil {
			for _, d := range devices {
				list.Store(uint32(d.SerialNumber), d)
			}
		}
	}()

	wg.Wait()

	keys := []uint32{}
	list.Range(func(key, value interface{}) bool {
		if _, ok := value.(types.Device); ok {
			keys = append(keys, key.(uint32))
		}

		return true
	})

	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	table := [][]string{}
	for _, key := range keys {
		if value, ok := list.Load(key); ok {
			device := value.(types.Device)
			record := []string{
				fmt.Sprintf("%v", device.Name),
				fmt.Sprintf("%v", uint32(device.SerialNumber)),
				fmt.Sprintf("%v", device.IpAddress.To4()),
				fmt.Sprintf("%v", device.SubnetMask.To4()),
				fmt.Sprintf("%v", device.Gateway.To4()),
				fmt.Sprintf("%v", device.MacAddress),
				fmt.Sprintf("%v", device.Version),
				fmt.Sprintf("%v", device.Date),
			}

			table = append(table, record)
		}
	}

	widths := []int{0, 0, 0, 0, 0, 0, 0, 0}
	for _, row := range table {
		for i, f := range row {
			if len(f) > widths[i] {
				widths[i] = len(f)
			}
		}
	}

	formats := []string{}
	if widths[0] > 0 {
		formats = append(formats, fmt.Sprintf("%%-%v[%v]s", widths[0], 1))
	}
	for i, w := range widths[1:] {
		formats = append(formats, fmt.Sprintf("%%-%v[%v]s", w, i+2))
	}

	format := strings.Join(formats, "  ")
	for _, row := range table {
		fmt.Printf(format, row[0], row[1], row[2], row[3], row[4], row[5], row[6], row[7])
		fmt.Println()
	}

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
