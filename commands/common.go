package commands

import (
	"errors"
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// WEEKDAYS
type days map[string]bool

func (d days) String() string {
	list := []string{}
	for _, v := range []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"} {
		if d[v] {
			list = append(list, v)
		}
	}

	return strings.Join(list, ",")
}

func (d days) parse(arg string) error {
	delete(d, "Monday")
	delete(d, "Tuesday")
	delete(d, "Wednesday")
	delete(d, "Thursday")
	delete(d, "Friday")
	delete(d, "Saturday")
	delete(d, "Sunday")

	tokens := strings.Split(arg, ",")
	for _, t := range tokens {
		day := strings.ToLower(t)
		switch {
		case strings.HasPrefix(day, "mon"):
			d["Monday"] = true
		case strings.HasPrefix(day, "tue"):
			d["Tuesday"] = true
		case strings.HasPrefix(day, "wed"):
			d["Wednesday"] = true
		case strings.HasPrefix(day, "thu"):
			d["Thursday"] = true
		case strings.HasPrefix(day, "fri"):
			d["Friday"] = true
		case strings.HasPrefix(day, "sat"):
			d["Saturday"] = true
		case strings.HasPrefix(day, "sun"):
			d["Sunday"] = true
		default:
			return fmt.Errorf("%v: unrecognised 'weekday'", t)
		}
	}

	return nil
}

func getUint8(index int, missing, invalid string) (uint8, error) {
	if len(flag.Args()) < index+1 {
		return 0, errors.New(missing)
	}

	valid, _ := regexp.MatchString("[0-9]+", flag.Arg(index))

	if !valid {
		return 0, fmt.Errorf(invalid, flag.Arg(index))
	}

	N, err := strconv.ParseUint(flag.Arg(index), 10, 8)

	if err != nil {
		return 0, fmt.Errorf(invalid, flag.Arg(index))
	}

	return uint8(N), err
}

func getUint32(index int, missing, invalid string) (uint32, error) {
	if len(flag.Args()) < index+1 {
		return 0, errors.New(missing)
	}

	valid, _ := regexp.MatchString("[0-9]+", flag.Arg(index))

	if !valid {
		return 0, fmt.Errorf(invalid, flag.Arg(index))
	}

	N, err := strconv.ParseUint(flag.Arg(index), 10, 32)

	if err != nil {
		return 0, fmt.Errorf(invalid, flag.Arg(index))
	}

	return uint32(N), err
}

func getString(index int, missing, invalid string) (string, error) {
	if len(flag.Args()) < index+1 {
		return "", errors.New(missing)
	}

	return flag.Arg(index), nil
}

func getDate(index int, missing, invalid string) (*time.Time, error) {
	if len(flag.Args()) < index+1 {
		return nil, errors.New(missing)
	}

	valid, _ := regexp.MatchString("[0-9]{4}-[0-9]{2}-[0-9]{2}", flag.Arg(index))

	if !valid {
		return nil, fmt.Errorf(invalid, flag.Arg(index))
	}

	date, err := time.Parse("2006-01-02", flag.Arg(index))

	if err != nil {
		return nil, fmt.Errorf(invalid, flag.Arg(index))
	}

	return &date, err
}

func getDoor(index int, missing, invalid string) (byte, error) {
	if len(flag.Args()) < index+1 {
		return 0, errors.New(missing)
	}

	valid, _ := regexp.MatchString("[1-4]", flag.Arg(index))

	if !valid {
		return 0, fmt.Errorf(invalid, flag.Arg(index))
	}

	door, err := strconv.Atoi(flag.Arg(index))

	if err != nil {
		return 0, fmt.Errorf(invalid, flag.Arg(index))
	}

	return byte(door), nil
}

func format(table [][]string) []string {
	rows := []string{}

	columns := 0
	for _, row := range table {
		if len(row) > columns {
			columns = len(row)
		}
	}

	widths := make([]int, columns)
	for _, row := range table {
		for i, f := range row {
			if len(f) > widths[i] {
				widths[i] = len(f)
			}
		}
	}

	formats := []string{}
	for _, w := range widths {
		formats = append(formats, fmt.Sprintf("%%-%vs", w))
	}

	for _, row := range table {
		line := []string{}
		for i, v := range row {
			line = append(line, fmt.Sprintf(formats[i], v))
		}

		rows = append(rows, strings.Join(line, "  "))
	}

	return rows
}
