package commands

import (
	"errors"
	"flag"
	"fmt"
	"net/netip"
	"strconv"
)

var SetListenerCmd = SetListener{}

type SetListener struct {
}

func (c *SetListener) Execute(ctx Context) error {
	if controller, err := getSerialNumber(ctx); err != nil {
		return err
	} else if addrport, err := c.getAddress(); err != nil {
		return err
	} else if interval, err := c.getInterval(); err != nil {
		return err
	} else if ok, err := ctx.uhppote.SetListener(controller, addrport, interval); err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("failed to set listener")
	} else {
		fmt.Printf("%-10v %v %d\n", controller, addrport, interval)

		return nil
	}
}

func (c *SetListener) CLI() string {
	return "set-listener"
}

func (c *SetListener) Description() string {
	return "Sets the IP address and port to which the controller sends access events"
}

func (c *SetListener) Usage() string {
	return "<serial number> <address:port> [interval]"
}

func (c *SetListener) Help() {
	fmt.Println("Usage: uhppote-cli [options] set-listener <serial number> <address:port> [interval]")
	fmt.Println()
	fmt.Println(" Sets the host address to which the controller sends access events")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println("  address        (required) IPv4 address")
	fmt.Println("  port           (required) IP port in the range 1 to 65535")
	fmt.Println("  interval       (optional) interval (in seconds) at which to send the most recent event")
	fmt.Println("                            (defaults to 0 - no auto-send)")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli set-listener 12345678  192.168.1.100:54321")
	fmt.Println("    uhppote-cli set-listener 12345678  192.168.1.100:54321 15")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *SetListener) RequiresConfig() bool {
	return false
}

func (c *SetListener) getAddress() (netip.AddrPort, error) {
	if len(flag.Args()) < 3 {
		return netip.AddrPort{}, errors.New("missing IPv4 address:port")
	} else if addr, err := netip.ParseAddrPort(flag.Arg(2)); err != nil {
		return netip.AddrPort{}, err
	} else if !addr.IsValid() {
		return netip.AddrPort{}, fmt.Errorf("invalid IPv4 address:port (%v)", flag.Arg(2))
	} else {
		return addr, nil
	}
}

func (c *SetListener) getInterval() (uint8, error) {
	if len(flag.Args()) < 4 {
		return 0, nil
	} else if interval, err := strconv.ParseUint(flag.Arg(3), 10, 8); err != nil {
		return 0, fmt.Errorf("invalid auto-send interval (%v)", flag.Arg(3))
	} else {
		return uint8(interval), nil
	}
}
