package commands

import (
	"bytes"
	"testing"

	"github.com/uhppoted/uhppote-core/types"
)

func TestGetCardPrint(t *testing.T) {
	getCard := GetCard{}
	expected := `10058400 2026-01-01 2026-12-31 Y N N 29 7531 1,4`

	card := types.Card{
		CardNumber: 10058400,
		From:       types.MustParseDate("2026-01-01"),
		To:         types.MustParseDate("2026-12-31"),
		Doors:      map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 29},
		PIN:        7531,
		FirstCard: types.FirstCardPrivileges{
			Door1: true,
			Door4: true,
		},
	}

	var b bytes.Buffer

	if err := getCard.print(card, &b); err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	} else if b.String() != expected {
		t.Errorf("Incorrectly formatted card\n   -- expected:\n%v\n   -- got:\n%v", expected, b.String())
	}
}

func TestGetCardPrintWithoutPINOrFirstCard(t *testing.T) {
	getCard := GetCard{}
	expected := `10058400 2026-01-01 2026-12-31 Y N N 29 - -`

	card := types.Card{
		CardNumber: 10058400,
		From:       types.MustParseDate("2026-01-01"),
		To:         types.MustParseDate("2026-12-31"),
		Doors:      map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 29},
	}

	var b bytes.Buffer

	if err := getCard.print(card, &b); err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	} else if b.String() != expected {
		t.Errorf("Incorrectly formatted card\n   -- expected:\n%v\n   -- got:\n%v", expected, b.String())
	}
}
