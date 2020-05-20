package commands

import (
	"github.com/uhppoted/uhppote-core/uhppote"
	"github.com/uhppoted/uhppoted-api/config"
)

// Context contains the environment and configuration information required for all commands
type Context struct {
	uhppote *uhppote.UHPPOTE
	config  *config.Config
}

// NewContext returns a valid Context initialized with the supplied UHPPOTE and
// configuration.
func NewContext(u *uhppote.UHPPOTE, c *config.Config) Context {
	return Context{u, c}
}

// Command defines the common functions for CLI command implementations. This will be
// replaced with the 'uhppoted-api' implementation in a future iteration.
type Command interface {
	Execute(context Context) error
	CLI() string
	Description() string
	Usage() string
	Help()
	RequiresConfig() bool
}
