package commands

import (
	"reflect"
	"testing"
)

func TestParseWeekdays(t *testing.T) {
	tests := []struct {
		arg  string
		days []string
	}{
		{arg: "Mon", days: []string{"Monday"}},
		{arg: "Tue", days: []string{"Tuesday"}},
		{arg: "Wed", days: []string{"Wednesday"}},
		{arg: "Thurs", days: []string{"Thursday"}},
		{arg: "Thu", days: []string{"Thursday"}},
		{arg: "Fri", days: []string{"Friday"}},
		{arg: "Sat", days: []string{"Saturday"}},
		{arg: "Sun", days: []string{"Sunday"}},

		{arg: "mon", days: []string{"Monday"}},
		{arg: "tue", days: []string{"Tuesday"}},
		{arg: "wed", days: []string{"Wednesday"}},
		{arg: "thurs", days: []string{"Thursday"}},
		{arg: "thu", days: []string{"Thursday"}},
		{arg: "fri", days: []string{"Friday"}},
		{arg: "sat", days: []string{"Saturday"}},
		{arg: "sun", days: []string{"Sunday"}},

		{arg: "Sat,Sun", days: []string{"Saturday", "Sunday"}},
	}

	for _, v := range tests {
		expected := days{}

		for _, d := range v.days {
			expected[d] = true
		}

		weekdays := days{
			"Monday":    true,
			"Tuesday":   true,
			"Wednesday": true,
			"Thursday":  true,
			"Friday":    true,
			"Saturday":  true,
			"Sunday":    true,
		}

		if err := weekdays.parse(v.arg); err != nil {
			t.Fatalf("Unexpected error parsing %v (%v)", v.arg, err)
		}

		if !reflect.DeepEqual(weekdays, expected) {
			t.Errorf("'%v' parsed incorrectly\n   expected: %v\n   got:      %v\n", v.arg, expected, weekdays)
		}
	}
}
