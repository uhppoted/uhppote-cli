package commands

import (
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	"os"
	"os/signal"
)

var ListenCmd = Listen{}

type Listen struct {
}

type listener struct {
}

func (l *listener) OnConnected() {
	fmt.Printf("Listening...\n")
}

func (l *listener) OnEvent(event *types.Status) {
	fmt.Printf("%v\n", event)
}

func (l *listener) OnError(err error) bool {
	fmt.Printf("ERROR: %v\n", err)
	return true
}

func (c *Listen) Execute(ctx Context) error {
	q := make(chan os.Signal)

	defer close(q)

	signal.Notify(q, os.Interrupt)

	return ctx.uhppote.Listen(&listener{}, q)
}

func (c *Listen) CLI() string {
	return "listen"
}

func (c *Listen) Description() string {
	return "Listens for access control events"
}

func (c *Listen) Usage() string {
	return ""
}

func (c *Listen) Help() {
	fmt.Println("Listens for access control events from UHPPOTE UT0311-L0x controllers configured to send events to this IP address and port")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *Listen) RequiresConfig() bool {
	return false
}
