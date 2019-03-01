package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"uhppote"
	"uhppote-cli/commands"
)

type bind struct {
	address net.UDPAddr
}

var VERSION = "v0.00.0"
var debug = false
var local = bind{net.UDPAddr{net.IPv4(0, 0, 0, 0), 60001, ""}}

func main() {
	flag.Var(&local, "bind", "Sets the local IP address and port to which to bind (e.g. 192.168.0.100:60001)")
	flag.BoolVar(&debug, "debug", false, "Displays vaguely useful information while processing a command")
	flag.Parse()

	u := uhppote.UHPPOTE{
		BindAddress: local.address,
		Debug:       debug,
	}

	command := parse()
	err := command.Execute(&u)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

func parse() commands.Command {
	var cmd commands.Command = commands.NewHelpCommand()
	var err error = nil

	if len(os.Args) > 1 {
		switch flag.Arg(0) {
		case "help":
			cmd = commands.NewHelpCommand()

		case "version":
			cmd = commands.NewVersionCommand(VERSION)

		case "list-devices":
			cmd, err = commands.NewListDevicesCommand()

		case "get-time":
			cmd, err = commands.NewGetTimeCommand()

		case "set-time":
			cmd, err = commands.NewSetTimeCommand()

		case "set-ip-address":
			cmd, err = commands.NewSetAddressCommand()

		case "get-authorised":
			cmd, err = commands.NewGetAuthorisedCommand()

		case "get-swipe":
			cmd, err = commands.NewGetSwipeCommand()

		case "authorise":
			cmd, err = commands.NewGrantCommand()
		}
	}

	if err == nil {
		return cmd
	}

	return commands.NewHelpCommand()
}

func (b *bind) String() string {
	return b.address.String()
}

func (b *bind) Set(s string) error {
	address, err := net.ResolveUDPAddr("udp", s)
	if err != nil {
		return err
	}

	b.address = *address

	return nil
}
