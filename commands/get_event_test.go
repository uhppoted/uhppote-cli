package commands

import (
	"testing"
)

func TestGetNextEventIndex(t *testing.T) {
	tests := []struct {
		first    uint32
		last     uint32
		current  uint32
		expected uint32
	}{
		// normal (last >= first)
		{first: 5, last: 37, current: 17, expected: 18},
		{first: 5, last: 37, current: 57, expected: 38},
		{first: 5, last: 37, current: 3, expected: 5},
	}

	cmd := GetEvent{}

	for _, v := range tests {
		next := cmd.getNextIndex(v.first, v.last, v.current)
		if next != v.expected {
			t.Errorf("Incorrect 'next' event - expected:%v, got:%v", v.expected, next)
		}
	}
}
