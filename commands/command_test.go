package commands

import (
	"fmt"
	"net"
	"net/netip"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	"github.com/uhppoted/uhppoted-lib/config"
)

func TestCommandNewContext(t *testing.T) {
	LA, _ := time.LoadLocation("America/Los_Angeles")

	u := stub{}
	c := config.Config{
		Devices: config.DeviceMap{
			405419896: &config.Device{
				Name:     "Alpha",
				Address:  types.MustParseControllerAddr("192.168.1.100:60000"),
				Doors:    []string{"Gryffindor", "Hufflepuff", "Ravenclaw", "Slytherin"},
				TimeZone: "America/Los_Angeles",
				Protocol: "tcp",
			},
			303986753: &config.Device{
				Name:    "Beta",
				Address: types.MustParseControllerAddr("192.168.1.100:60000"),
				Doors:   []string{"Great Hall", "Kitchen", "Dungeon", "Hogsmeade"},
			},
		},
	}

	expected := Context{
		uhppote: &u,
		config:  &c,
		devices: []uhppote.Device{
			uhppote.Device{
				Name:     "Beta",
				DeviceID: 303986753,
				Doors:    []string{"Great Hall", "Kitchen", "Dungeon", "Hogsmeade"},
				Address:  types.MustParseControllerAddr("192.168.1.100:60000"),
				TimeZone: time.Local,
				Protocol: "udp",
			},
			uhppote.Device{
				Name:     "Alpha",
				DeviceID: 405419896,
				Doors:    []string{"Gryffindor", "Hufflepuff", "Ravenclaw", "Slytherin"},
				Address:  types.MustParseControllerAddr("192.168.1.100:60000"),
				TimeZone: LA,
				Protocol: "tcp",
			},
		},
		debug: false,
	}

	ctx := NewContext(&u, &c, false)

	if !reflect.DeepEqual(ctx, expected) {
		t.Errorf("incorrect context\n   expected:%v\n   got:     %v", expected, ctx)
	}
}

type stub struct {
}

func (s *stub) GetDevices() ([]types.Device, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *stub) GetDevice(controller uint32) (*types.Device, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *stub) SetAddress(controller uint32, address, mask, gateway net.IP) (*types.Result, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *stub) GetTime(controller uint32) (*types.Time, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *stub) SetTime(controller uint32, datetime time.Time) (*types.Time, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *stub) GetDoorControlState(controller uint32, door byte) (*types.DoorControlState, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *stub) SetDoorControlState(controller uint32, door uint8, state types.ControlState, delay uint8) (*types.DoorControlState, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *stub) SetDoorPasscodes(controller uint32, door uint8, passcodes ...uint32) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (s *stub) SetInterlock(controller uint32, interlock types.Interlock) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (s *stub) ActivateKeypads(controller uint32, keypads map[uint8]bool) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (s *stub) GetAntiPassback(controller uint32) (types.AntiPassback, error) {
	return types.Disabled, fmt.Errorf("Not implemented")
}

func (s *stub) SetAntiPassback(controller uint32, antipassback types.AntiPassback) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (s *stub) GetListener(controller uint32) (netip.AddrPort, uint8, error) {
	return netip.AddrPort{}, 0, fmt.Errorf("Not implemented")
}

func (s *stub) SetListener(controller uint32, address netip.AddrPort, interval uint8) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (s *stub) GetStatus(controller uint32) (*types.Status, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *stub) GetCards(controller uint32) (uint32, error) {
	return 0, fmt.Errorf("Not implemented")
}

func (s *stub) GetCardByIndex(controller, index uint32) (*types.Card, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *stub) GetCardByID(controller, cardID uint32) (*types.Card, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *stub) PutCard(controller uint32, card types.Card, formats ...types.CardFormat) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (s *stub) DeleteCard(controller uint32, cardNumber uint32) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (s *stub) DeleteCards(controller uint32) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (s *stub) GetTimeProfile(controller uint32, profileID uint8) (*types.TimeProfile, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *stub) SetTimeProfile(controller uint32, profile types.TimeProfile) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (s *stub) ClearTimeProfiles(controller uint32) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (s *stub) ClearTaskList(controller uint32) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (s *stub) AddTask(controller uint32, task types.Task) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (s *stub) RefreshTaskList(controller uint32) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (s *stub) RecordSpecialEvents(controller uint32, enable bool) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (s *stub) GetEvent(controller, index uint32) (*types.Event, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *stub) GetEventIndex(controller uint32) (*types.EventIndex, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *stub) SetEventIndex(controller, index uint32) (*types.EventIndexResult, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *stub) Listen(listener uhppote.Listener, q chan os.Signal) error {
	return fmt.Errorf("Not implemented")
}

func (s *stub) OpenDoor(controller uint32, door uint8) (*types.Result, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *stub) SetPCControl(controller uint32, enable bool) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (s *stub) RestoreDefaultParameters(controller uint32) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (s *stub) DeviceList() map[uint32]uhppote.Device {
	return nil
}

func (s *stub) ListenAddrList() []netip.AddrPort {
	return nil
}
