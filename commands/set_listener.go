package commands

import (
	"errors"
	"flag"
	"fmt"
	"net"
)

var SetListenerCmd = SetListener{}

type SetListener struct {
}

func (c *SetListener) Execute(ctx Context) error {
	serialNumber, err := getSerialNumber(ctx)
	if err != nil {
		return err
	}

	if len(flag.Args()) < 3 {
		return errors.New("missing IP address")
	}

	address, err := net.ResolveUDPAddr("udp", flag.Arg(2))
	if err != nil {
		return err
	}

	if address == nil || address.IP.To4() == nil {
		return fmt.Errorf("invalid UDP address: %v", flag.Arg(2))
	}

	listener, err := ctx.uhppote.SetListener(serialNumber, *address)

	if err == nil {
		fmt.Printf("%v\n", listener)
	}

	return err
}

func (c *SetListener) CLI() string {
	return "set-listener"
}

func (c *SetListener) Description() string {
	return "Sets the IP address and port to which the controller sends access events"
}

func (c *SetListener) Usage() string {
	return "<serial number> <address:port>"
}

func (c *SetListener) Help() {
	fmt.Println("Usage: uhppote-cli [options] set-listener <serial number> <address:port>")
	fmt.Println()
	fmt.Println(" Sets the host address to which the controller sends access events")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println("  address        (required) IPv4 address")
	fmt.Println("  port           (required) IP port in the range 1 to 65535")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli set-listener 12345678  192.168.1.100:54321")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *SetListener) RequiresConfig() bool {
	return false
}
