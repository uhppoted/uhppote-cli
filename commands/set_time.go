package commands

import (
	"flag"
	"fmt"
	"time"
)

var SetTimeCmd = SetTime{}

type SetTime struct {
}

func (c *SetTime) Execute(ctx Context) error {
	serialNumber, err := getSerialNumber(ctx)
	if err != nil {
		return err
	}

	datetime := time.Now()
	if len(flag.Args()) > 2 {
		if flag.Arg(2) == "now" {
			datetime = time.Now()
		} else {
			datetime, err = time.Parse("2006-01-02 15:04:05", flag.Arg(2))
			if err != nil {
				return fmt.Errorf("invalid date/time parameter: %v", flag.Arg(3))
			}
		}
	}

	devicetime, err := ctx.uhppote.SetTime(serialNumber, datetime)

	if err == nil {
		fmt.Printf("%s\n", devicetime)
	}

	return err
}

func (c *SetTime) CLI() string {
	return "set-time"
}

func (c *SetTime) Description() string {
	return "Sets the controller internal clock"
}

func (c *SetTime) Usage() string {
	return "<serial number> [now|<yyyy-mm-dd HH:mm:ss>]"
}

func (c *SetTime) Help() {
	fmt.Println("Usage: uhppote-cli [options] set-time <serial number> [command options]")
	fmt.Println()
	fmt.Println(" Sets the controller date/time to the supplied time. Defaults to 'now'. Command format")
	fmt.Println()
	fmt.Println(" <serial number> [now|<yyyy-mm-dd HH:mm:ss>]")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    -debug  Displays internal information for diagnosing errors")
	fmt.Println()
	fmt.Println("  Command options:")
	fmt.Println()
	fmt.Println("    now                    Sets the controller time to the system time of the local system")
	fmt.Println("    'yyyy-mm-dd HH:mm:ss'  Sets the controller time to the explicitly supplied instant")
	fmt.Println()
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli set-time")
	fmt.Println("    uhppote-cli set-time now")
	fmt.Println("    uhppote-cli set-time '2019-01-12 20:15:32'")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *SetTime) RequiresConfig() bool {
	return false
}
