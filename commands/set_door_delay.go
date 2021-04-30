package commands

import (
	"fmt"
)

var SetDoorDelayCmd = SetDoorDelay{}

type SetDoorDelay struct {
}

func (c *SetDoorDelay) Execute(ctx Context) error {
	serialNumber, err := getSerialNumber(ctx)
	if err != nil {
		return err
	}

	door, err := getDoor(2, "Missing door", "Invalid door: %v")
	if err != nil {
		return err
	}

	delay, err := getUint8(3, "Missing delay", "Invalid delay: %v")
	if err != nil {
		return err
	}

	state, err := ctx.uhppote.GetDoorControlState(serialNumber, door)
	if err != nil {
		return err
	}

	record, err := ctx.uhppote.SetDoorControlState(serialNumber, door, state.ControlState, delay)
	if err != nil {
		return err
	}

	fmt.Printf("%s %v %v\n", record.SerialNumber, record.Door, record.Delay)

	return nil
}

func (c *SetDoorDelay) CLI() string {
	return "set-door-delay"
}

func (c *SetDoorDelay) Description() string {
	return "Sets the duration for which a door lock is kept open"
}

func (c *SetDoorDelay) Usage() string {
	return "<serial number> <door> <delay>"
}

func (c *SetDoorDelay) Help() {
	fmt.Println("Usage: uhppote-cli [options] set-door-delay <serial number> <door> <delay>")
	fmt.Println()
	fmt.Println(" Sets the door open delay (in seconds), independently of the door control state")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println("  door           (required) door (1,2,3 or 4")
	fmt.Println("  delay          (required) delay in seconds")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli set-door-delay 12345678 3 15")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *SetDoorDelay) RequiresConfig() bool {
	return false
}
