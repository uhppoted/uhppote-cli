package commands

import (
	"slices"
	"testing"
)

func TestPutCardWithFirstCard(t *testing.T) {
	putCard := PutCard{}

	tests := []struct {
		arg      string
		expected []uint8
	}{
		{"", []uint8{}},
		{"-", []uint8{}},
		{"1", []uint8{1}},
		{"1,2", []uint8{1, 2}},
		{"1,2,2", []uint8{1, 2}},
		{"1,2,3", []uint8{1, 2, 3}},
		{"1,2,3,4", []uint8{1, 2, 3, 4}},
		{"1,2,3,3,4", []uint8{1, 2, 3, 4}},
		{"4,1,3,2", []uint8{1, 2, 3, 4}},
		{"1,2,3,X", []uint8{1, 2, 3}},
	}

	for _, test := range tests {
		if firstcard, err := putCard.getFirstCard(test.arg); err != nil {
			t.Fatalf("Unexpected error (%v)", err)
		} else if !slices.Equal(firstcard, test.expected) {
			t.Errorf("Incorrect first-card privileges - expected:%v\n, got:%v", test.expected, firstcard)
		}
	}
}
