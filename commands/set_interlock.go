package commands

import (
	"flag"
	"fmt"
	"regexp"

	"github.com/uhppoted/uhppote-core/types"
)

var SetInterlockCmd = SetInterlock{}

type SetInterlock struct {
}

func (c *SetInterlock) Execute(ctx Context) error {
	controllerID, err := getSerialNumber(ctx)
	if err != nil {
		return err
	}

	if interlock, err := c.getInterlock(); err != nil {
		return err
	} else if ok, err := ctx.uhppote.SetInterlock(controllerID, interlock); err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("%v  failed to set interlock %v", controllerID, interlock)
	} else {
		fmt.Printf("%v  set interlock %v\n", controllerID, interlock)
	}

	return nil
}

func (c *SetInterlock) CLI() string {
	return "set-interlock"
}

func (c *SetInterlock) Description() string {
	return "Sets the door interlock"
}

func (c *SetInterlock) Usage() string {
	return "<serial number> <interlock>"
}

func (c *SetInterlock) Help() {
	fmt.Println("Usage: uhppote-cli [options] set-interlock <serial number> <interlock>")
	fmt.Println()
	fmt.Println(" Sets the door interlock")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println("  interlock      (required) interlock")
	fmt.Println("                            - none")
	fmt.Println("                            - 1&2,3&4")
	fmt.Println("                            - 1&3,2&4")
	fmt.Println("                            - 1&2&3")
	fmt.Println("                            - 1&2&3&4")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli set-interlocl 405419896 1&2&3")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *SetInterlock) RequiresConfig() bool {
	return false
}

func (c *SetInterlock) getInterlock() (types.Interlock, error) {
	if len(flag.Args()) < 3 {
		return types.NoInterlock, fmt.Errorf("missing interlock")
	}

	interlock := flag.Arg(2)

	switch {
	case regexp.MustCompile("none").MatchString(interlock):
		return types.NoInterlock, nil

	case regexp.MustCompile("1&2,3&4").MatchString(interlock):
		return types.Interlock12_34, nil

	case regexp.MustCompile("1&3,2&4").MatchString(interlock):
		return types.Interlock13_24, nil

	case regexp.MustCompile("1&2&3").MatchString(interlock):
		return types.Interlock123, nil

	case regexp.MustCompile("1&2&3&4").MatchString(interlock):
		return types.Interlock1234, nil

	default:
		return types.NoInterlock, fmt.Errorf("invalid interlock (%v)", interlock)
	}
}
