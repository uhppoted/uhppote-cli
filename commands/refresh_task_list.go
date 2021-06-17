package commands

import (
	"fmt"

	"github.com/uhppoted/uhppoted-lib/config"
)

var RefreshTaskListCmd = RefreshTaskList{}

type RefreshTaskList struct {
}

func (c *RefreshTaskList) Execute(ctx Context) error {
	serialNumber, err := getSerialNumber(ctx)
	if err != nil {
		return err
	}

	refreshed, err := ctx.uhppote.RefreshTaskList(serialNumber)
	if err != nil {
		return err
	}

	fmt.Printf("%v %v\n", serialNumber, refreshed)

	return nil
}

func (c *RefreshTaskList) CLI() string {
	return "refresh-task-list"
}

func (c *RefreshTaskList) Description() string {
	return "Refreshes updated task list on the controller"
}

func (c *RefreshTaskList) Usage() string {
	return "<serial number>"
}

func (c *RefreshTaskList) Help() {
	fmt.Println("Usage: uhppote-cli [options] refresh-task-list <serial number>")
	fmt.Println()
	fmt.Println(" Refreshes an updated task list on a controller")
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
	fmt.Println("    uhppote-cli --debug refresh-task-list 9876543210")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *RefreshTaskList) RequiresConfig() bool {
	return false
}
