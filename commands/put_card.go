package commands

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-lib/config"
)

var PutCardCmd = PutCard{}

type PutCard struct {
}

func (c *PutCard) Execute(ctx Context) error {
	serialNumber, err := getSerialNumber(ctx)
	if err != nil {
		return err
	}

	cardNumber, err := getUint32(2, "missing card number", "invalid card number: %v")
	if err != nil {
		return err
	}

	from, err := getDate(3, "missing start date", "invalid start date: %v")
	if err != nil {
		return err
	}

	to, err := getDate(4, "missing end date", "invalid end date: %v")
	if err != nil {
		return err
	}

	permissions, err := getPermissions()
	if err != nil {
		return err
	}

	for _, door := range []uint8{1, 2, 3, 4} {
		if v, ok := permissions[door]; ok && v >= 2 && v <= 254 {
			if profile, err := ctx.uhppote.GetTimeProfile(serialNumber, uint8(v)); err != nil {
				return err
			} else if profile == nil {
				return fmt.Errorf("time profile %v is not defined", v)
			}
		}
	}

	pin, err := getPIN()
	if err != nil {
		return err
	}

	start := types.Date(*from)
	end := types.Date(*to)
	authorised, err := ctx.uhppote.PutCard(serialNumber, types.Card{
		CardNumber: cardNumber,
		From:       &start,
		To:         &end,
		Doors:      permissions,
		PIN:        pin,
	})

	if err != nil {
		return err
	}

	fmt.Printf("%v %v %v\n", serialNumber, cardNumber, authorised)

	return nil
}

func (c *PutCard) CLI() string {
	return "put-card"
}

func (c *PutCard) Description() string {
	return "Sets the card access rights on a single access controller"
}

func (c *PutCard) Usage() string {
	return "<serial number> <card number> <start date> <end date> <doors>"
}

func (c *PutCard) Help() {
	fmt.Println("Usage: uhppote-cli [options] put-card <serial number> <card number> <start date> <end date> <doors>")
	fmt.Println()
	fmt.Println(" Adds (or updates) a card to the list of the cards managed by a controller")
	fmt.Println()
	fmt.Println("  <serial number>  (required) controller serial number")
	fmt.Println("  <card number>    (required) card number")
	fmt.Println("  <start date>     (required) start date YYYY-MM-DD")
	fmt.Println("  <end date>       (required) end date   YYYY-MM-DD")
	fmt.Println("  <doors>          (required) list of permitted doors [1,2,3,4]. Unlisted doors will be set to 'N'")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --config  File path for the 'conf' file containing the controller configuration")
	fmt.Printf("              (defaults to %s)\n", config.DefaultConfig)
	fmt.Println("    --debug   Displays internal information for diagnosing errors")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli put-card 12345678 918273645 2020-01-01 2020-12-31 1,2,4")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *PutCard) RequiresConfig() bool {
	return false
}

func getPermissions() (map[uint8]uint8, error) {
	index := 5
	permissions := map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}

	if len(flag.Args()) > index {
		tokens := strings.Split(flag.Arg(index), ",")

		for _, token := range tokens {
			match := regexp.MustCompile("([1-4])(?::([0-9]+))?").FindStringSubmatch(token)
			if match == nil || len(match) < 3 {
				return nil, fmt.Errorf("invalid door '%v'", token)
			}

			door, err := strconv.ParseInt(match[1], 10, 8)
			if err != nil {
				return nil, fmt.Errorf("invalid door ID '%v' (%v)", match[1], err)
			}

			if match[2] == "" {
				permissions[uint8(door)] = 1
			} else {
				profile, err := strconv.ParseUint(match[2], 10, 8)
				if err != nil {
					return nil, fmt.Errorf("invalid time profile '%v' (%v)", match[2], err)
				} else if profile < 2 || profile > 254 {
					return nil, fmt.Errorf("invalid time profile '%v' (valid profiles are in the range 2 to 254)", match[2])
				}

				permissions[uint8(door)] = uint8(profile)
			}
		}
	}

	return permissions, nil
}

func getPIN() (types.PIN, error) {
	pin := types.PIN(0)

	if len(flag.Args()) > 6 {
		arg := flag.Arg(6)
		if ok := regexp.MustCompile("[0-9]+").MatchString(arg); !ok {
			return pin, fmt.Errorf("invalid PIN (%v)", arg)
		} else if v, err := strconv.ParseUint(arg, 10, 32); err != nil {
			return pin, err
		} else if v > 999999 {
			return pin, fmt.Errorf("invalid PIN (%v)", v)
		} else {
			pin = types.PIN(v)
		}
	}

	return pin, nil
}
