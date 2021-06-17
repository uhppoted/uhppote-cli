package commands

import (
	"fmt"

	"github.com/uhppoted/uhppoted-lib/config"
)

var ClearTaskListCmd = ClearTaskList{}

type ClearTaskList struct {
}

func (c *ClearTaskList) Execute(ctx Context) error {
	serialNumber, err := getSerialNumber(ctx)
	if err != nil {
		return err
	}

	cleared, err := ctx.uhppote.ClearTaskList(serialNumber)
	if err != nil {
		return err
	}

	fmt.Printf("%v %v\n", serialNumber, cleared)

	return nil
}

func (c *ClearTaskList) CLI() string {
	return "clear-task-list"
}

func (c *ClearTaskList) Description() string {
	return "Clears all tasks from the controller"
}

func (c *ClearTaskList) Usage() string {
	return "<serial number>"
}

func (c *ClearTaskList) Help() {
	fmt.Println("Usage: uhppote-cli [options] clear-task-list <serial number>")
	fmt.Println()
	fmt.Println(" Clears all tasks from the controller")
	fmt.Println()
	fmt.Println("  <serial number>  (required) controller serial number")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --config  File path for the 'conf' file containing the controller configuration")
	fmt.Printf("              (defaults to %s)\n", config.DefaultConfig)
	fmt.Println("    --debug   Displays internal information for diagnosing errors")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli --debug clear-task-list 9876543210")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *ClearTaskList) RequiresConfig() bool {
	return false
}
