module github.com/uhppoted/uhppote-cli

go 1.14

require (
	github.com/uhppoted/uhppote-core v0.6.2
	github.com/uhppoted/uhppoted-api v0.6.2
)

replace (
    github.com/uhppoted/uhppote-core => ../uhppote-core
    github.com/uhppoted/uhppoted-api => ../uhppoted-api
)
