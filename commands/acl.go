package commands

import (
	"github.com/uhppoted/uhppote-core/uhppote"
	"sort"
)

func getDevices(ctx *Context) []*uhppote.Device {
	keys := []uint32{}
	for id, _ := range ctx.config.Devices {
		keys = append(keys, id)
	}

	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	devices := []*uhppote.Device{}
	for _, id := range keys {
		d := ctx.config.Devices[id]
		devices = append(devices, uhppote.NewDevice(id, d.Address, d.Rollover, d.Doors))
	}

	return devices
}
