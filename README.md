# uhppote-cli

CLI for the *UHPPOTE UT0311-L0x* TCP/IP Wiegand access control boards. 

Supported operating systems:
- Linux
- MacOS
- Windows

## Raison d'être

The manufacturer supplied application is 'Windows-only' and provides limited support for integration with other
systems.

## Releases

- v0.5.1: Initial release following restructuring into standalone Go *modules* and *git submodules*

## Installation

### Building from source

#### Dependencies

| *Dependency*                          | *Description*                                          |
| ------------------------------------- | ------------------------------------------------------ |
| [com.github/uhppoted/uhppote-core][1] | Device level API implementation                        |
| golang.org/x/lint/golint              | Additional *lint* check for release builds             |

### Binaries

## uhppote-cli

Usage: *uhppote-cli [--bind <address:port>] [--debug] \<command\> \<arguments\>*

Supported commands:

- help
- version
- get-devices
- get-device
- set-address
- get-status
- get-time
- set-time
- get-door-delay
- set-door-delay
- get-door-control
- set-door-control
- get-listener
- set-listener
- get-cards
- get-card
- get-events
- get-swipe-index
- set-event-index
- open
- grant
- revoke
- revoke-all
- load-acl
- listen


[1]: https://github.com/uhppoted/uhppote-core
