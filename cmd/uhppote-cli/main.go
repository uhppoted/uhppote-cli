package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/uhppoted/uhppote-cli/commands"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	"github.com/uhppoted/uhppoted-lib/config"
)

var cli = []commands.Command{
	&commands.VersionCmd,
	&commands.GetDevicesCmd,
	&commands.GetDeviceCmd,
	&commands.SetAddressCmd,
	&commands.GetListenerCmd,
	&commands.SetListenerCmd,
	&commands.GetTimeCmd,
	&commands.SetTimeCmd,
	&commands.GetDoorDelayCmd,
	&commands.SetDoorDelayCmd,
	&commands.GetDoorControlCmd,
	&commands.SetDoorControlCmd,
	&commands.RecordSpecialEventsCmd,
	&commands.GetStatusCmd,
	&commands.GetCardsCmd,
	&commands.GetCardCmd,
	&commands.PutCardCmd,
	&commands.DeleteCardCmd,
	&commands.DeleteCardsCmd,
	&commands.GetTimeProfileCmd,
	&commands.GetTimeProfilesCmd,
	&commands.SetTimeProfileCmd,
	&commands.SetTimeProfilesCmd,
	&commands.ClearTimeProfilesCmd,
	&commands.ClearTaskListCmd,
	&commands.RefreshTaskListCmd,
	&commands.AddTaskCmd,
	&commands.SetTaskListCmd,
	&commands.ShowCmd,
	&commands.GrantCmd,
	&commands.RevokeCmd,
	&commands.LoadACLCmd,
	&commands.GetACLCmd,
	&commands.CompareACLCmd,
	&commands.GetEventsCmd,
	&commands.GetEventCmd,
	&commands.GetEventIndexCmd,
	&commands.SetEventIndexCmd,
	&commands.OpenDoorCmd,
	&commands.SetPCControlCmd,
	&commands.SetInterlockCmd,
	&commands.ActivateKeypadsCmd,
	&commands.ListenCmd,
}

var options = struct {
	config    string
	bind      types.BindAddr
	broadcast types.BroadcastAddr
	listen    types.ListenAddr
	timeout   time.Duration
	debug     bool
}{}

func main() {
	// ... parse command line args
	var bind types.BindAddr
	var broadcast types.BroadcastAddr
	var listen types.ListenAddr

	flag.StringVar(&options.config, "config", options.config, "Specifies the path for the config file")
	flag.Var(&bind, "bind", "Sets the local IP address and port to which to bind (e.g. 192.168.0.100:60001)")
	flag.Var(&broadcast, "broadcast", "Sets the IP address and port for UDP broadcast (e.g. 192.168.0.255:60000)")
	flag.Var(&listen, "listen", "Sets the local IP address and port to which to bind for events (e.g. 192.168.0.100:60001)")
	flag.DurationVar(&options.timeout, "timeout", 2500*time.Millisecond, "Sets the timeout for a response from a controller (e.g. 3.5s)")
	flag.BoolVar(&options.debug, "debug", options.debug, "Displays internal information for diagnosing errors")
	flag.Parse()

	cmd, err := parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n   ERROR: %v\n\n", err)
		os.Exit(1)
	}

	if cmd == nil {
		help()
		return
	}

	// initialise execution context
	conf := configuration(cmd)

	if conf.BindAddress != nil {
		options.bind = *conf.BindAddress
	}

	if conf.BroadcastAddress != nil {
		options.broadcast = *conf.BroadcastAddress
	}

	if conf.ListenAddress != nil {
		options.listen = *conf.ListenAddress
	}

	devices := []uhppote.Device{}
	for s, d := range conf.Devices {
		// ... because d is *Device and all devices end up with the same info if you don't make a manual copy
		name := d.Name
		deviceID := s
		address := d.Address
		doors := d.Doors

		if device := uhppote.NewDevice(name, deviceID, address, doors); device != nil {
			devices = append(devices, *device)
		}
	}

	// ... override defaults/conf settings with command line options
	overrides := func(a *flag.Flag) {
		switch a.Name {
		case "bind":
			options.bind = bind

		case "broadcast":
			options.broadcast = broadcast

		case "listen":
			options.listen = listen
		}
	}

	flag.Visit(overrides)

	if err := validate(bind, broadcast, listen); err != nil {
		fmt.Fprintf(os.Stderr, "\n   ERROR: %v\n\n", err)
		os.Exit(1)
	}

	u := uhppote.NewUHPPOTE(options.bind, options.broadcast, options.listen, options.timeout, devices, options.debug)

	// execute command
	ctx := commands.NewContext(u, conf, options.debug)
	err = cmd.Execute(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n   ERROR: %v\n\n", err)
		os.Exit(1)
	}
}

func validate(bind types.BindAddr, broadcast types.BroadcastAddr, listen types.ListenAddr) error {
	// validate bind.address port
	port := bind.Port

	if port != 0 && port == broadcast.Port {
		return fmt.Errorf("bind address port (%v) must not be the same as the broadcast address port", bind.Port)
	}

	if port != 0 && port == listen.Port {
		return fmt.Errorf("bind address port (%v) must not be the same as the listen address port", bind.Port)
	}

	return nil
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
			fmt.Fprintf(os.Stderr, "\n   ERROR: %v\n\n", err)
			os.Exit(1)
		}
	} else {
		info, err := os.Stat(config.DefaultConfig)
		if err != nil {
			if !os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "\n   WARN:  %v\n\n", err)
			} else if cmd.RequiresConfig() {
				fmt.Fprintf(os.Stderr, "\n   ERROR: '%s' requires a valid configuration file:\n", cmd.CLI())
				fmt.Fprintf(os.Stderr, "          %v\n\n", err)
				os.Exit(1)
			}
		} else if !info.IsDir() {
			if err := conf.Load(config.DefaultConfig); err != nil {
				if cmd.RequiresConfig() {
					fmt.Fprintf(os.Stderr, "\n   ERROR: '%s' requires a valid configuration file:\n", cmd.CLI())
					fmt.Fprintf(os.Stderr, "          %v\n\n", err)
					os.Exit(1)
				} else {
					fmt.Fprintf(os.Stderr, "\n   WARN:  %v\n", err)
				}
			} else if options.debug || cmd.RequiresConfig() {
				fmt.Fprintf(os.Stderr, "\n ... using default configuration from %v\n", config.DefaultConfig)
			}
		}
	}

	if err := conf.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "\n   ERROR: %v\n\n", err)
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

			fmt.Fprintf(os.Stderr, "Invalid command: %v. Type 'help commands' to get a list of supported commands\n", flag.Arg(1))
			return
		}
	}

	usage()
}

func usage() {
	format := "    %-21s %s\n"
	fmt.Println()
	fmt.Println("  Usage: uhppote-cli [options] <command>")
	fmt.Println()
	fmt.Println("  Commands:")
	fmt.Println()
	fmt.Printf(format, "help", "Displays this message")
	fmt.Printf(format, "", "For help on a specific command use 'uhppote-cli help <command>'")
	fmt.Println()

	for _, c := range cli {
		fmt.Printf(format, c.CLI(), c.Description())
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
		fmt.Printf(" %-21s %s\n", c.CLI(), c.Usage())
	}

	fmt.Println()
}
