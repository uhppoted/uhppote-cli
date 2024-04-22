package commands

import (
	"flag"
	"fmt"
	"net"
	"net/netip"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/uhppoted/uhppote-core/uhppote"
	"github.com/uhppoted/uhppoted-lib/config"
)

// Context contains the environment and configuration information required for all commands
type Context struct {
	uhppote uhppote.IUHPPOTE
	devices []uhppote.Device
	config  *config.Config
	debug   bool
}

// NewContext returns a valid Context initialized with the supplied UHPPOTE and
// configuration.
func NewContext(u uhppote.IUHPPOTE, c *config.Config, debug bool) Context {
	keys := []uint32{}
	for id := range c.Devices {
		keys = append(keys, id)
	}

	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	devices := []uhppote.Device{}
	for _, id := range keys {
		d := c.Devices[id]
		address := resolve(d.Address)

		if device := uhppote.NewDevice(d.Name, id, address, d.Doors); device != nil {
			devices = append(devices, *device)
		}
	}

	return Context{
		uhppote: u,
		devices: devices,
		config:  c,
		debug:   debug,
	}
}

// Command defines the common functions for CLI command implementations. This will be
// replaced with the 'uhppoted-lib' implementation in a future iteration.
type Command interface {
	Execute(context Context) error
	CLI() string
	Description() string
	Usage() string
	Help()
	RequiresConfig() bool
}

func clean(s string) string {
	return regexp.MustCompile(`[\s\t]+`).ReplaceAllString(strings.ToLower(s), "")
}

func getSerialNumber(ctx Context) (uint32, error) {
	return getSerialNumberI(ctx, 1)
}

func getSerialNumberI(ctx Context, index int) (uint32, error) {
	if len(flag.Args()) < index+1 {
		return 0, fmt.Errorf("missing controller serial number")
	}

	arg := flag.Arg(index)

	// lookup controller by name
	if ctx.config != nil {
		for k, v := range ctx.config.Devices {
			if clean(arg) == clean(v.Name) {
				return k, nil
			}
		}
	}

	// numeric serial number?
	if valid, _ := regexp.MatchString("[0-9]+", arg); !valid {
		return 0, fmt.Errorf("invalid controller serial number:%v", arg)
	}

	if N, err := strconv.ParseUint(arg, 10, 32); err != nil {
		return 0, fmt.Errorf("invalid controller serial number (%v)", arg)
	} else {
		return uint32(N), nil
	}
}

func resolve(udp *net.UDPAddr) *netip.AddrPort {
	if udp == nil {
		return nil
	}

	addr := netip.AddrFrom4([4]byte(udp.IP.To4()))
	port := uint16(udp.Port)
	address := netip.AddrPortFrom(addr, port)

	return &address
}
