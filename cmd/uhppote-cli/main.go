package main

import (
	"flag"
	"fmt"
	"github.com/uhppoted/uhppote-cli/commands"
	"github.com/uhppoted/uhppote-core/uhppote"
	"github.com/uhppoted/uhppoted-api/config"
	"net"
	"os"
)

type addr struct {
	address *net.UDPAddr
}

var cli = []commands.Command{
	&commands.VersionCmd,
	&commands.GetDevicesCmd,
	&commands.GetDeviceCmd,
	&commands.SetAddressCmd,
	&commands.GetStatusCmd,
	&commands.GetTimeCmd,
	&commands.SetTimeCmd,
	&commands.GetDoorDelayCmd,
	&commands.SetDoorDelayCmd,
	&commands.GetDoorControlCmd,
	&commands.SetDoorControlCmd,
	&commands.GetListenerCmd,
	&commands.SetListenerCmd,
	&commands.GetCardsCmd,
	&commands.GetCardCmd,
	&commands.PutCardCmd,
	&commands.DeleteCardCmd,
	&commands.DeleteAllCmd,
	&commands.ShowCmd,
	&commands.GrantCmd,
	&commands.RevokeCmd,
	&commands.LoadACLCmd,
	&commands.GetACLCmd,
	&commands.CompareACLCmd,
	&commands.GetEventsCmd,
	&commands.GetEventIndexCmd,
	&commands.SetEventIndexCmd,
	&commands.OpenDoorCmd,
	&commands.ListenCmd,
}

var options = struct {
	config    string
	bind      addr
	broadcast addr
	listen    addr
	debug     bool
}{
	config:    "",
	bind:      addr{nil},
	broadcast: addr{nil},
	listen:    addr{nil},
	debug:     false,
}

func main() {
	flag.StringVar(&options.config, "config", options.config, "Specifies the path for the config file")
	flag.Var(&options.bind, "bind", "Sets the local IP address and port to which to bind (e.g. 192.168.0.100:60001)")
	flag.Var(&options.broadcast, "broadcast", "Sets the IP address and port for UDP broadcast (e.g. 192.168.0.255:60000)")
	flag.Var(&options.listen, "listen", "Sets the local IP address and port to which to bind for events (e.g. 192.168.0.100:60001)")
	flag.BoolVar(&options.debug, "debug", options.debug, "Displays internal information for diagnosing errors")
	flag.Parse()

	cmd, err := parse()
	if err != nil {
		fmt.Printf("\n   ERROR: %v\n\n", err)
		os.Exit(1)
	}

	if cmd == nil {
		help()
		return
	}

	// initialise execution context
	conf := configuration(cmd)
	u := uhppote.UHPPOTE{
		Devices: make(map[uint32]*uhppote.Device),
		Debug:   options.debug,
	}

	u.BindAddress = conf.BindAddress
	u.BroadcastAddress = conf.BroadcastAddress
	u.ListenAddress = conf.ListenAddress

	for s, d := range conf.Devices {
		u.Devices[s] = uhppote.NewDevice(s, d.Address, d.Rollover, d.Doors)
	}

	if options.bind.address != nil {
		u.BindAddress = options.bind.address
		u.ListenAddress = options.bind.address
	}

	if options.broadcast.address != nil {
		u.BroadcastAddress = options.broadcast.address
	}

	if options.listen.address != nil {
		u.ListenAddress = options.listen.address
	}

	ctx := commands.NewContext(&u, conf)

	// execute command
	err = cmd.Execute(ctx)
	if err != nil {
		fmt.Printf("\n   ERROR: %v\n\n", err)
		os.Exit(1)
	}
}

// Optionally loads the configuration from file, falling back to the default configuration file
// if a file is not specified by the --conf command line option. For 'device' and 'miscellaneous'
// commands the configuration file is optional and a note is posted in debug mode if the
// default configuration file is being used. A valid configuration file is mandatory for ACL
// commands - a note is posted in debug mode if the default configuration file is in use,
// but this is probably the desired behaviour.
func configuration(cmd commands.Command) *config.Config {
	conf := config.NewConfig()

	if options.config != "" {
		if err := conf.Load(options.config); err != nil {
			fmt.Printf("\n   ERROR: %v\n\n", err)
			os.Exit(1)
		}
	} else {
		info, err := os.Stat(config.DefaultConfig)
		if err != nil {
			if !os.IsNotExist(err) {
				fmt.Printf("\n   WARN:  %v\n\n", err)
			} else if cmd.RequiresConfig() {
				fmt.Printf("\n   ERROR: '%s' requires a valid configuration file:\n", cmd.CLI())
				fmt.Printf("          %v\n\n", err)
				os.Exit(1)
			}
		} else if !info.IsDir() {
			if err := conf.Load(config.DefaultConfig); err != nil {
				if cmd.RequiresConfig() {
					fmt.Printf("\n   ERROR: '%s' requires a valid configuration file:\n", cmd.CLI())
					fmt.Printf("          %v\n\n", err)
					os.Exit(1)
				} else {
					fmt.Printf("\n   WARN:  %v\n", err)
				}
			} else if options.debug || cmd.RequiresConfig() {
				fmt.Printf("\n ... using default configuration from %v\n", config.DefaultConfig)
			}
		}
	}

	if err := conf.Validate(); err != nil {
		fmt.Printf("\n   ERROR: %v\n\n", err)
		os.Exit(1)
	}

	return conf
}

func parse() (commands.Command, error) {
	var cmd commands.Command = nil
	var err error = nil

	if len(os.Args) > 1 {
		for _, c := range cli {
			if c.CLI() == flag.Arg(0) {
				cmd = c
			}
		}
	}

	return cmd, err
}

func (b *addr) String() string {
	return b.address.String()
}

func (b *addr) Set(s string) error {
	address, err := net.ResolveUDPAddr("udp", s)
	if err != nil {
		return err
	}

	b.address = address

	return nil
}

func help() {
	if len(flag.Args()) > 0 && flag.Arg(0) == "help" {
		if len(flag.Args()) > 1 {

			if flag.Arg(1) == "commands" {
				helpCommands()
				return
			}

			for _, c := range cli {
				if c.CLI() == flag.Arg(1) {
					c.Help()
					return
				}
			}

			fmt.Printf("Invalid command: %v. Type 'help commands' to get a list of supported commands\n", flag.Arg(1))
			return
		}
	}

	usage()
}

func usage() {
	fmt.Println()
	fmt.Println("  Usage: uhppote-cli [options] <command>")
	fmt.Println()
	fmt.Println("  Commands:")
	fmt.Println()
	fmt.Println("    help             Displays this message")
	fmt.Println("                     For help on a specific command use 'uhppote-cli help <command>'")

	for _, c := range cli {
		fmt.Printf("    %-16s %s\n", c.CLI(), c.Description())
	}

	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --config    Sets the configuration file")
	fmt.Println("    --bind      Sets the local IP address and port to use")
	fmt.Println("    --broadcast Sets the IP address and port to use for UDP broadcast")
	fmt.Println("    --listen    Sets the local IP address and port to use for receiving device events")
	fmt.Println("    --debug     Displays internal information for diagnosing errors")
	fmt.Println()
}

func helpCommands() {
	fmt.Println("Supported commands:")
	fmt.Println()

	for _, c := range cli {
		fmt.Printf(" %-16s %s\n", c.CLI(), c.Usage())
	}

	fmt.Println()
}
