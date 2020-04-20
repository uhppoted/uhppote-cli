module github.com/uhppoted/uhppote-cli

go 1.14

require (
	github.com/uhppoted/uhppote-core v0.6.1-0.20200410194200-582e323e671e
	github.com/uhppoted/uhppoted-api v0.6.1-0.20200410194243-1e2448e29254
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/tools v0.0.0-20200420001825-978e26b7c37c // indirect
)

replace (
	github.com/uhppoted/uhppote-core => ../uhppote-core
	github.com/uhppoted/uhppoted-api => ../uhppoted-api
)
