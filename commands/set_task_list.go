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

var SetTaskListCmd = SetTaskList{}

type SetTaskList struct {
}

func (c *SetTaskList) Execute(ctx Context) error {
	serialNumber, err := getSerialNumber(ctx)
	if err != nil {
		return err
	}

	file, err := c.getTSVFile()
	if err != nil {
		return err
	} else if file == "" {
		return fmt.Errorf("missing TSV file with tasks")
	}

	tasks, err := c.parse(file)
	if err != nil {
		return err
	} else if tasks == nil {
		return fmt.Errorf("could not extract tasks from TSV File '%s'", file)
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

func (c *SetTaskList) CLI() string {
	return "set-task-list"
}

func (c *SetTaskList) Description() string {
	return "Sets the list of tasks on a controller"
}

func (c *SetTaskList) Usage() string {
	return "<serial number>"
}

func (c *SetTaskList) Help() {
	fmt.Println("Usage: uhppote-cli [options] set-task-list <serial number> <file>")
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
	fmt.Println("    uhppote-cli set-task-list 9876543210 tasks.tsv")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *SetTaskList) RequiresConfig() bool {
	return false
}

func (c *SetTaskList) getTSVFile() (string, error) {
	if len(flag.Args()) < 3 {
		return "", nil
	}

	file := flag.Arg(2)
	stat, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return file, fmt.Errorf("file '%s' does not exist", file)
		} else {
			return "", err
		}
	}

	if stat.Mode().IsDir() {
		return "", fmt.Errorf("file '%s' is a directory", file)
	}

	if !stat.Mode().IsRegular() {
		return "", fmt.Errorf("file '%s' is not a real file", file)
	}

	return file, nil
}

func (c *SetTaskList) load(ctx Context, serialNumber uint32, tasks []types.Task) (int, []error, error) {
	warnings := []error{}
	created := [][]string{}

	for id, task := range tasks {
		if err := c.validate(task); err != nil {
			warnings = append(warnings, fmt.Errorf("task %-3v %v", id+1, err))
			continue
		}

		if ok, err := ctx.uhppote.AddTask(serialNumber, task); err != nil {
			return len(created), nil, err
		} else if !ok {
			warnings = append(warnings, fmt.Errorf("%v: could not create task definition %v", serialNumber, id+1))
		} else {
			cards := ""
			if task.Task == types.EnableMoreCards {
				cards = fmt.Sprintf("%d", task.Cards)
			}

			created = append(created, []string{
				fmt.Sprintf("%v", id+1),
				fmt.Sprintf("%v", task.Task),
				fmt.Sprintf("%v", task.Door),
				fmt.Sprintf("%v:%v", task.From, task.To),
				fmt.Sprintf("%v", task.Weekdays),
				fmt.Sprintf("%v", task.Start),
				cards,
			})
		}
	}

	rows := format(created)
	for _, v := range rows {
		fmt.Printf("   ... created task definition %s\n", v)
	}

	return len(created), warnings, nil
}

func (c *SetTaskList) validate(task types.Task) error {
	if task.To.Before(task.From) {
		return fmt.Errorf("'To' date (%v) is before 'From' date (%v)", task.To, task.From)
	}

	return nil
}

func (c *SetTaskList) parse(file string) ([]types.Task, error) {
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
		task := types.Task{
			Task: record.Task,
			Door: record.Door,
			From: record.From,
			To:   record.To,
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
