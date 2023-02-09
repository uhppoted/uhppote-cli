package commands

import (
	"fmt"

	"github.com/uhppoted/uhppote-core/types"
)

var SetDoorControlCmd = SetDoorControl{}

type SetDoorControl struct {
}

func (c *SetDoorControl) Execute(ctx Context) error {
	states := map[string]types.ControlState{
		"normally open":   1,
		"normally closed": 2,
		"controlled":      3,
	}

	serialNumber, err := getSerialNumber(ctx)
	if err != nil {
		return err
	}

	door, err := getDoor(2, "missing door", "invalid door: %v")
	if err != nil {
		return err
	}

	control, err := getString(3, "missing control value", "invalid control value: %v")
	if err != nil {
		return err
	} else if _, ok := states[control]; !ok {
		return fmt.Errorf("invalid door control value: %s (expected 'normally open', 'normally closed' or 'controlled'", control)
	}

	state, err := ctx.uhppote.GetDoorControlState(serialNumber, door)
	if err != nil {
		return err
	}

	record, err := ctx.uhppote.SetDoorControlState(serialNumber, door, states[control], state.Delay)
	if err != nil {
		return err
	}

	fmt.Printf("%s %v %v (%v)\n", record.SerialNumber, record.Door, uint8(record.ControlState), record.ControlState)

	return nil
}

func (c *SetDoorControl) CLI() string {
	return "set-door-control"
}

func (c *SetDoorControl) Description() string {
	return "Sets the control state (normally open, normally close or controlled) for a door"
}

func (c *SetDoorControl) Usage() string {
	return "<serial number> <door> <state>"
}

func (c *SetDoorControl) Help() {
	fmt.Println("Usage: uhppote-cli [options] set-door-control <serial number> <door> <state>")
	fmt.Println()
	fmt.Println(" Sets the door control state")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println("  door           (required) door (1,2,3 or 4")
	fmt.Println("  state          (required) 'normally open','normally closed', 'controlled'")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli set-door-control 12345678 3 'normally open'")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *SetDoorControl) RequiresConfig() bool {
	return false
}
