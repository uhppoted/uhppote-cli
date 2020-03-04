package parsers

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/uhppoted/uhppote-cli/config"
	"github.com/uhppoted/uhppote-core/types"
	"io"
	"strconv"
	"strings"
	"time"
)

type ACL map[uint32]map[uint32]*types.Card

type index struct {
	cardnumber int
	from       int
	to         int
	doors      map[uint32][]int
}

func (a *ACL) Load(f *bufio.Reader, path string, cfg *config.Config) (*ACL, error) {
	acl := make(ACL)
	for id, _ := range cfg.Devices {
		acl[id] = make(map[uint32]*types.Card)
	}

	r := csv.NewReader(f)
	r.Comma = '\t'

	header, err := r.Read()
	if err != nil {
		return nil, err
	}

	index, err := parseHeader(header, path, cfg)
	if err != nil {
		return nil, err
	}
	line := 0

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		line += 1

		cards, err := parseRecord(record, index, path)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Line %d: %v\n", line, err))
		}

		for id, card := range *cards {
			if acl[id] != nil {
				if acl[id][card.CardNumber] != nil {
					return nil, errors.New(fmt.Sprintf("Duplicate card number (%v)\n", card.CardNumber))
				}

				acl[id][card.CardNumber] = &card
			}
		}
	}

	return &acl, nil
}

func parseHeader(header []string, path string, cfg *config.Config) (*index, error) {
	columns := make(map[string]int)

	index := index{
		cardnumber: 0,
		from:       0,
		to:         0,
		doors:      make(map[uint32][]int),
	}

	for id, _ := range cfg.Devices {
		index.doors[id] = make([]int, 4)
	}

	for c, field := range header {
		key := strings.ReplaceAll(strings.ToLower(field), " ", "")
		ix := c + 1

		if columns[key] != 0 {
			return nil, errors.New(fmt.Sprintf("Duplicate column name '%s' in File '%s", field, path))
		}

		columns[key] = ix
	}

	index.cardnumber = columns["cardnumber"]
	index.from = columns["from"]
	index.to = columns["to"]

	for id, device := range cfg.Devices {
		for i, door := range device.Door {
			if d := strings.ReplaceAll(strings.ToLower(door), " ", ""); d != "" {
				index.doors[id][i] = columns[d]
			}
		}
	}

	if index.cardnumber == 0 {
		return nil, errors.New(fmt.Sprintf("File '%s' does not include a column 'Card Number'", path))
	}

	if index.from == 0 {
		return nil, errors.New(fmt.Sprintf("File '%s' does not include a column 'From'", path))
	}

	if index.to == 0 {
		return nil, errors.New(fmt.Sprintf("File '%s' does not include a column 'to'", path))
	}

	for id, device := range cfg.Devices {
		for i, door := range device.Door {
			if d := strings.ReplaceAll(strings.ToLower(door), " ", ""); d != "" {
				if index.doors[id][i] == 0 {
					return nil, errors.New(fmt.Sprintf("File '%s' does not include a column for door '%s'", path, door))
				}
			}
		}
	}

	return &index, nil
}

func parseRecord(record []string, index *index, path string) (*map[uint32]types.Card, error) {
	cards := make(map[uint32]types.Card, 0)

	for k, v := range index.doors {
		cardno, err := getCardNumber(record, index)
		if err != nil {
			return nil, err
		}

		from, err := getFromDate(record, index)
		if err != nil {
			return nil, err
		}

		to, err := getToDate(record, index)
		if err != nil {
			return nil, err
		}

		doors, err := getDoors(record, v)
		if err != nil {
			return nil, err
		}

		cards[k] = types.Card{
			CardNumber: cardno,
			From:       *from,
			To:         *to,
			Doors:      doors,
		}
	}

	return &cards, nil
}

func getCardNumber(record []string, index *index) (uint32, error) {
	f := get(record, index.cardnumber)
	cardnumber, err := strconv.ParseUint(f, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("Invalid card number '%s' (%w)", f, err)
	}

	return uint32(cardnumber), nil
}

func getFromDate(record []string, index *index) (*types.Date, error) {
	f := get(record, index.from)
	date, err := time.ParseInLocation("2006-01-02", f, time.Local)
	if err != nil {
		return nil, fmt.Errorf("Invalid 'from' date '%s' (%w)", f, err)
	}

	from := types.Date(date)

	return &from, nil
}

func getToDate(record []string, index *index) (*types.Date, error) {
	f := get(record, index.to)
	date, err := time.ParseInLocation("2006-01-02", f, time.Local)
	if err != nil {
		return nil, fmt.Errorf("Invalid 'to' date '%s' (%w)", f, err)
	}

	to := types.Date(date)

	return &to, nil
}

func getDoors(record []string, v []int) ([]bool, error) {
	doors := make([]bool, 4)

	for i, d := range v {
		if d == 0 {
			doors[i] = false
			continue
		}

		switch get(record, d) {
		case "Y":
			doors[i] = true
		case "N":
			doors[i] = false
		default:
			return doors, fmt.Errorf("Expected 'Y/N' for door: '%s'", record[d])
		}
	}

	return doors, nil
}

func get(record []string, ix int) string {
	return strings.TrimSpace(record[ix-1])
}
