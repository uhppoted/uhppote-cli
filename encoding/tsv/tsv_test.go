package tsv

import (
	"reflect"
	"testing"

	"github.com/uhppoted/uhppote-core/types"
)

var data = []byte(`Profile	From	To	Mon	Tue	Wed	Thurs	Fri	Sat	Sun	Start1	End1	Start2	End2	Start3	End3	Linked
2	2021-04-01	2021-12-31	N	N	Y	N	Y	N	Y	08:30	11:30	00:00	00:00	13:45	17:00	
3	2021-04-01	         	Y	Y	Y	N	N	N	N	08:35	11:30	00:00	13:15	17:05  	17:15	2
29	2021-04-01	2021-12-31	Y	N	Y	Y	N	Y	N	09:45	11:30	00:00	    	13:45	     	
`)

type profile struct {
	ID        int         `tsv:"Profile"`
	From      types.Date  `tsv:"From"`
	To        *types.Date `tsv:"To"`
	Monday    bool        `tsv:"Mon"`
	Tuesday   bool        `tsv:"Tue"`
	Wednesday bool        `tsv:"Wed"`
	Thursday  bool        `tsv:"Thurs"`
	Friday    bool        `tsv:"Fri"`
	Saturday  bool        `tsv:"Sat"`
	Sunday    bool        `tsv:"Sun"`
	Start1    types.HHmm  `tsv:"Start1"`
	End1      *types.HHmm `tsv:"End1"`
	Start2    types.HHmm  `tsv:"Start2"`
	End2      *types.HHmm `tsv:"End2"`
	Start3    types.HHmm  `tsv:"Start3"`
	End3      *types.HHmm `tsv:"End3"`
	Linked    int         `tsv:"Linked"`
}

func TestTSVUnmarshal(t *testing.T) {
	expected := []profile{
		profile{
			ID:        2,
			From:      date("2021-04-01"),
			To:        pdate("2021-12-31"),
			Monday:    false,
			Tuesday:   false,
			Wednesday: true,
			Thursday:  false,
			Friday:    true,
			Saturday:  false,
			Sunday:    true,
			Start1:    hhmm("08:30"),
			End1:      phhmm("11:30"),
			Start2:    hhmm("00:00"),
			End2:      phhmm("00:00"),
			Start3:    hhmm("13:45"),
			End3:      phhmm("17:00"),
			Linked:    0,
		},

		profile{
			ID:        3,
			From:      date("2021-04-01"),
			Monday:    true,
			Tuesday:   true,
			Wednesday: true,
			Thursday:  false,
			Friday:    false,
			Saturday:  false,
			Sunday:    false,
			Start1:    hhmm("08:35"),
			End1:      phhmm("11:30"),
			Start2:    hhmm("00:00"),
			End2:      phhmm("13:15"),
			Start3:    hhmm("17:05"),
			End3:      phhmm("17:15"),
			Linked:    2,
		},

		profile{
			ID:        29,
			From:      date("2021-04-01"),
			To:        pdate("2021-12-31"),
			Monday:    true,
			Tuesday:   false,
			Wednesday: true,
			Thursday:  true,
			Friday:    false,
			Saturday:  true,
			Sunday:    false,
			Start1:    hhmm("09:45"),
			End1:      phhmm("11:30"),
			Start2:    hhmm("00:00"),
			Start3:    hhmm("13:45"),
			Linked:    0,
		},
	}

	profiles := []profile{}

	if err := Unmarshal(data, &profiles); err != nil {
		t.Fatalf("Unexpected error unmarshalling TSV (%v)", err)
	}

	if !reflect.DeepEqual(profiles, expected) {
		if len(profiles) != len(expected) {
			t.Errorf("Expected %v profiles, got:%v", len(expected), len(profiles))
		} else {
			for i, row := range expected {
				if !reflect.DeepEqual(profiles[i], row) {
					t.Errorf("Row %v\n   expected:%v\n   got:     %v", i+1, row, profiles[i])
				}
			}
		}
	}
}

func date(s string) types.Date {
	d, _ := types.DateFromString(s)

	return *d
}

func pdate(s string) *types.Date {
	d, _ := types.DateFromString(s)

	return d
}

func hhmm(s string) types.HHmm {
	t, _ := types.HHmmFromString(s)

	return *t
}

func phhmm(s string) *types.HHmm {
	t, _ := types.HHmmFromString(s)

	return t
}
