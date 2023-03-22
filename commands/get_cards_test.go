package commands

import (
	"bytes"
	"testing"
	"time"

	"github.com/uhppoted/uhppote-core/types"
)

func TestGetCardsPrint(t *testing.T) {
	getCards := GetCards{}
	expected := `12345    2023-01-01 2023-12-21 Y N N N
8165539  2023-01-01 2023-12-31 Y N N 29
8165538  2023-01-01 2023-12-31 Y N N 29 7531
`

	recordset := []types.Card{
		types.Card{
			CardNumber: 12345,
			From:       types.ToDate(2023, time.January, 1),
			To:         types.ToDate(2023, time.December, 21),
			Doors:      map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0},
		},
		types.Card{
			CardNumber: 8165539,
			From:       types.ToDate(2023, time.January, 1),
			To:         types.ToDate(2023, time.December, 31),
			Doors:      map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 29},
		},
		types.Card{
			CardNumber: 8165538,
			From:       types.ToDate(2023, time.January, 1),
			To:         types.ToDate(2023, time.December, 31),
			Doors:      map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 29},
			PIN:        7531,
		},
	}

	var b bytes.Buffer

	if err := getCards.print(recordset, &b); err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	} else if b.String() != expected {
		t.Errorf("Incorrectly formatted cards\n   -- expected:\n%v\n   -- got:\n%v", expected, b.String())
	}
}

func TestGetCardsPrintWithInvalidCardNumber(t *testing.T) {
	getCards := GetCards{}
	expected := `12345     2023-01-01 2023-12-21 Y N N N
8165539   2023-01-01 2023-12-31 Y N N 29
8165538   2023-01-01 2023-12-31 Y N N 29 7531
192837465 2023-01-01 2023-12-31 Y N N 29 7531
`

	recordset := []types.Card{
		types.Card{
			CardNumber: 12345,
			From:       types.ToDate(2023, time.January, 1),
			To:         types.ToDate(2023, time.December, 21),
			Doors:      map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0},
		},
		types.Card{
			CardNumber: 8165539,
			From:       types.ToDate(2023, time.January, 1),
			To:         types.ToDate(2023, time.December, 31),
			Doors:      map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 29},
		},
		types.Card{
			CardNumber: 8165538,
			From:       types.ToDate(2023, time.January, 1),
			To:         types.ToDate(2023, time.December, 31),
			Doors:      map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 29},
			PIN:        7531,
		},
		types.Card{
			CardNumber: 192837465,
			From:       types.ToDate(2023, time.January, 1),
			To:         types.ToDate(2023, time.December, 31),
			Doors:      map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 29},
			PIN:        7531,
		},
	}

	var b bytes.Buffer

	if err := getCards.print(recordset, &b); err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	} else if b.String() != expected {
		t.Errorf("Incorrectly formatted cards\n   -- expected:\n%v\n   -- got:\n%v", expected, b.String())
		// t.Errorf("Incorrectly formatted cards\n   -- expected:\n%v\n   -- got:\n%v", []byte(expected), b.Bytes())
	}
}
