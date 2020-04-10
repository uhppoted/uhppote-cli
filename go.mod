module github.com/uhppoted/uhppote-cli

go 1.14

require (
	github.com/uhppoted/uhppote-core v0.6.1-0.20200410194200-582e323e671e
	github.com/uhppoted/uhppoted-api v0.6.1-0.20200410194243-1e2448e29254
)

replace (
	github.com/uhppoted/uhppote-core => ../uhppote-core
	github.com/uhppoted/uhppoted-api => ../uhppoted-api
)

