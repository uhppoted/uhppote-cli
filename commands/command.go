package commands

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/uhppoted/uhppote-core/uhppote"
	"github.com/uhppoted/uhppoted-api/config"
)

// Context contains the environment and configuration information required for all commands
type Context struct {
	uhppote *uhppote.UHPPOTE
	config  *config.Config
}

// NewContext returns a valid Context initialized with the supplied UHPPOTE and
// configuration.
func NewContext(u *uhppote.UHPPOTE, c *config.Config) Context {
	return Context{u, c}
}

// Command defines the common functions for CLI command implementations. This will be
// replaced with the 'uhppoted-api' implementation in a future iteration.
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
	if len(flag.Args()) < 2 {
		return 0, fmt.Errorf("Missing controller serial number")
	}

	arg := flag.Arg(1)

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
		return 0, fmt.Errorf("Invalid controller serial number:%v", arg)
	}

	if N, err := strconv.ParseUint(arg, 10, 32); err != nil {
		return 0, fmt.Errorf("Invalid controller serial number (%v)", arg)
	} else {
		return uint32(N), nil
	}
}
