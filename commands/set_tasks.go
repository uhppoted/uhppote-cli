package commands

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-lib/config"
	"github.com/uhppoted/uhppoted-lib/encoding/tsv"
)

var SetTasksCmd = SetTasks{}

type SetTasks struct {
}

func (c *SetTasks) Execute(ctx Context) error {
	serialNumber, err := getSerialNumber(ctx)
	if err != nil {
		return err
	}

	file, err := c.getTSVFile()
	if err != nil {
		return err
	} else if file == "" {
		return fmt.Errorf("Missing TSV file with tasks")
	}

	tasks, err := c.parse(file)
	if err != nil {
		return err
	} else if tasks == nil {
		return fmt.Errorf("Could not extract tasks from TSV File '%s'", file)
	} else if len(tasks) == 0 {
		fmt.Printf("   WARNING File '%s' does not contain any valid task definitions\n", file)
	}

	// ... clear task list
	cleared, err := ctx.uhppote.ClearTaskList(serialNumber)
	if err != nil {
		return err
	} else if !cleared {
		return fmt.Errorf("could not clear task list")
	}

	fmt.Printf("   ... %v cleared task list\n", serialNumber)

	// ... set tasks
	created, warnings, err := c.load(ctx, serialNumber, tasks)
	if err != nil {
		return err
	}

	if created == 0 {
		fmt.Printf("   ... %v created %v tasks\n", serialNumber, created)
	}

	// ... refresh task list
	refreshed, err := ctx.uhppote.RefreshTaskList(serialNumber)
	if err != nil {
		return err
	} else if !refreshed {
		return fmt.Errorf("could not refresh task list")
	}

	fmt.Printf("   ... %v refreshed task list\n", serialNumber)

	// ... done
	if len(warnings) > 0 {
		fmt.Println()
		for _, warning := range warnings {
			fmt.Printf("   WARN  %v\n", warning)
		}
		fmt.Println()
	}

	return nil
}

func (c *SetTasks) CLI() string {
	return "set-tasks"
}

func (c *SetTasks) Description() string {
	return "Sets the list of tasks on a controller"
}

func (c *SetTasks) Usage() string {
	return "<serial number>"
}

func (c *SetTasks) Help() {
	fmt.Println("Usage: uhppote-cli [options] set-tasks <serial number> <file>")
	fmt.Println()
	fmt.Println(" Clears any existing task defined on a controller, adds the task defined in the file and then invokes")
	fmt.Println(" refresh-tasks to activate the new task list.")
	fmt.Println()
	fmt.Println("  serial number  (required) controller serial number")
	fmt.Println("  file           (required) TSV file with list of tasks")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --config  File path for the 'conf' file containing the controller configuration")
	fmt.Printf("              (defaults to %s)\n", config.DefaultConfig)
	fmt.Println("    --debug   Displays internal information for diagnosing errors")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli set-tasks 9876543210 tasks.tsv")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *SetTasks) RequiresConfig() bool {
	return false
}

func (c *SetTasks) getTSVFile() (string, error) {
	if len(flag.Args()) < 3 {
		return "", nil
	}

	file := flag.Arg(2)
	stat, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return file, fmt.Errorf("File '%s' does not exist", file)
		} else {
			return "", err
		}
	}

	if stat.Mode().IsDir() {
		return "", fmt.Errorf("File '%s' is a directory", file)
	}

	if !stat.Mode().IsRegular() {
		return "", fmt.Errorf("File '%s' is not a real file", file)
	}

	return file, nil
}

func (c *SetTasks) load(ctx Context, serialNumber uint32, tasks []types.Task) (int, []error, error) {
	warnings := []error{}
	count := 0

	for id, task := range tasks {
		if err := c.validate(task); err != nil {
			warnings = append(warnings, fmt.Errorf("task %-3v %v", id+1, err))
			continue
		}

		if ok, err := ctx.uhppote.AddTask(serialNumber, task); err != nil {
			return count, nil, err
		} else if !ok {
			warnings = append(warnings, fmt.Errorf("%v: could not create task definition %v", serialNumber, id+1))
		} else {
			fmt.Printf("   ... created task defintion %v\n", id+1)
			count++
		}
	}

	return count, warnings, nil
}

func (c *SetTasks) validate(task types.Task) error {
	if task.From == nil {
		return fmt.Errorf("invalid 'From' date (%v)", task.From)
	}

	if task.To == nil {
		return fmt.Errorf("invalid 'To' date (%v)", task.To)
	}

	if task.To.Before(*task.From) {
		return fmt.Errorf("'To' date (%v) is before 'From' date (%v)", task.To, task.From)
	}

	return nil
}

func (c *SetTasks) parse(file string) ([]types.Task, error) {
	type tsvTask struct {
		Task      types.TaskType `tsv:"Task"`
		Door      uint8          `tsv:"Door"`
		From      types.Date     `tsv:"From"`
		To        types.Date     `tsv:"To"`
		Monday    bool           `tsv:"Mon"`
		Tuesday   bool           `tsv:"Tue"`
		Wednesday bool           `tsv:"Wed"`
		Thursday  bool           `tsv:"Thurs"`
		Friday    bool           `tsv:"Fri"`
		Saturday  bool           `tsv:"Sat"`
		Sunday    bool           `tsv:"Sun"`
		Start     types.HHmm     `tsv:"Start"`
		Cards     uint8          `tsv:"Cards"`
	}

	recordset := []tsvTask{}

	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if err := tsv.Unmarshal(bytes, &recordset); err != nil {
		return nil, err
	}

	tasks := []types.Task{}
	for _, record := range recordset {
		from := record.From
		to := record.To

		task := types.Task{
			Task: record.Task,
			Door: record.Door,
			From: &from,
			To:   &to,
			Weekdays: types.Weekdays{
				time.Monday:    record.Monday,
				time.Tuesday:   record.Tuesday,
				time.Wednesday: record.Wednesday,
				time.Thursday:  record.Thursday,
				time.Friday:    record.Friday,
				time.Saturday:  record.Saturday,
				time.Sunday:    record.Sunday,
			},
			Start: record.Start,
			Cards: record.Cards,
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
