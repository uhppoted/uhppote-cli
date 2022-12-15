![build](https://github.com/uhppoted/uhppote-cli/workflows/build/badge.svg)

# uhppote-cli

CLI for the *UHPPOTE UT0311-L0x* TCP/IP Wiegand access control boards. 

Supported operating systems:
- Linux
- MacOS
- Windows
- Raspberry Pi (ARM7)

## Raison d'être

The CLI provides a cross-platform set of command-line functions that allow direct interaction with a UHPPOTE access controller,
as well as facilitating scripting for automation of routine tasks e.g. synchronizing controller system time across daylight 
savings changes.

## Releases

| *Version* | *Description*                                                                                      |
| --------- | -------------------------------------------------------------------------------------------------- |
| v0.8.3    | Added ARM64 to release builds                                                                      |
| v0.8.2    | Maintenance release for version compatibility with `uhppote-core` v0.8.2                           |
| v0.8.1    | Maintenance release for version compatibility with `uhppote-core` v0.8.1                           |
| v0.8.0    | Maintenance release for version compatibility with `uhppote-core` v0.8.0                           |
| v0.7.3    | Maintenance release for version compatibility with `uhppote-core` v0.7.3                           |
| v0.7.2    | Replaced event rollover with `overwritten` event handling                                          |
| v0.7.1    | Added task list commands                                                                           |
| v0.7.0    | Added commands to get-, set- and clear time profiles                                               |
| v0.6.12   | Added validation for `bind`, `broadcast` and `listen` ports                                        |
| v0.6.10   | Maintenance release for version compatibility with `uhppoted-app-wild-apricot`                     |
| v0.6.8    | Maintenance release for version compatibility with `uhppote-core` `v0.6.8`                         |
| v0.6.7    | Implements `record-special-events` command to enable/disable door events                           |
| v0.6.5    | Maintenance release for version compatibility with `node-red-contrib-uhppoted`                     |
| v0.6.4    | Maintenance release for version compatibility with `uhppoted-app-sheets`                           |
| v0.6.3    | Reworked get-cards to handle deleted records                                                       |
| v0.6.2    | Fixed get-events for controllers without any events and improved configuration filse handling      |
| v0.6.1    | Added ACL commands to simplify managing card permissions across multiple controllers               |
| v0.6.0    | Maintenance release to keep compatibility with updated `uhppote-core`                              |
| v0.5.1    | Initial release following restructuring into standalone Go *modules* and *git submodules*          |

## Installation

An archive containing executables for all the supported operating systems can be downloaded from the [releases](https://github.com/uhppoted/uhppote-cli/releases) page. Alternatively, operating system specific tarballs can be found in the [uhppoted](https://github.com/uhppoted/uhppoted/releases) releases.

Installation is straightforward - download the archive and extract it to a directory of your choice and place the executable in a directory in your PATH. 

### `uhppoted.conf`

By default, `uhppote-cli` uses the communal `uhppoted.conf` configuration file shared by all the `uhppoted` project modules which is (_or will eventually be_) documented in [uhppoted](https://github.com/uhppoted/uhppoted). 

A valid configuration file is required **only** for system configurations where controllers are not findable on the local LAN (i.e. cannot receive and reply to UDP broadcasts) or for use with the _ACL_ commands which do require a valid `uhppoted.conf` to resolve the door ID to controller + door number. For _device_ commands the configuration file will used if present otherwise the internal default configuration will be used.

An alternative configuration file can be specified with the `--config` command line option, e.g.:

```
uhppote-cli --config ./uhppote.mine get-device 4726234734
```

A sample [uhppoted.conf](https://github.com/uhppoted/uhppoted/blob/master/runtime/simulation/405419896.conf) file is included in the `uhppoted` distribution.

### Building from source

Assuming you have `Go` and `make` installed:

```
git clone https://github.com/uhppoted/uhppote-cli.git
cd uhppote-cli
make build
```

If you prefer not to use `make`:
```
git clone https://github.com/uhppoted/uhppote-cli.git
cd uhppote-cli
mkdir bin
go build -trimpath -o bin ./...
```

The above commands build the `uhppote-cli` executable to the `bin` directory.

#### Dependencies

| *Dependency*                                                                 | *Description*                              |
| ---------------------------------------------------------------------------- | ------------------------------------------ |
| [com.github/uhppoted/uhppote-core](https://github.com/uhppoted/uhppote-core) | Device level API implementation            |
| [com.github/uhppoted/uhppoted-lib](https://github.com/uhppoted/uhppoted-lib) | common API for external applications       |


## uhppote-cli

Usage: ```uhppote-cli [options] <command> <parameters>```

General_ commands:

- `help`
- `version`

Device commands:

- [`get-devices`](#get-devices)
- [`get-device`](#get-device)
- [`set-address`](#set-address)
- [`get-time`](#get-time)
- [`set-time`](#set-time)
- [`get-door-delay`](#get-door-delay)
- [`set-door-delay`](#set-door-delay)
- [`get-door-control`](#get-door-control)
- [`set-door-control`](#set-door-control)
- [`get-listener`](#get-listener)
- [`set-listener`](#set-listener)
- [`record-special-events`](#record-special-events)
- [`get-status`](#get-status)
- [`get-cards`](#get-cards)
- [`get-card`](#get-card)
- [`put-card`](#put-card)
- [`delete-card`](#delete-card)
- [`delete-all`](#delete-all)
- [`get-time-profile`](#get-time-profile)
- [`set-time-profile`](#set-time-profile)
- [`get-time-profiles`](#get-time-profiles)
- [`set-time-profiles`](#set-time-profiles)
- [`clear-time-profiles`](#clear-time-profiles)
- [`clear-task-list`](#clear-task-list)
- [`add-task`](#add-task)
- [`refresh-task-list`](#refresh-task-list)
- [`set-task-list`](#set-task-list)
- [`get-events`](#get-events)
- [`get-event`](#get-event)
- [`get-event-index`](#get-event-index)
- [`set-event-index`](#set-event-index)
- [`open`](#open)
- [`listen`](#listen)

ACL commands:

- `grant`
- `revoke`
- `load-acl`
- `get-acl`
- `compare-acl`

#### Command options:
```
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Example:
  
   ./uhppote-cli --debug --config ./uhppoted.local get-time 4156216363
```

### General

#### `help`

Displays the usage information and a list of available commands. Command specific help displays the detailed usage for that command.

```
uhppote-cli help

  Examples:

  uhppote-cli help
  uhppote-cli help set-time
```

#### `version`

Displays the current application version.

```
uhppote-cli version

  Example:

  uhppote-cli version
```

### Device commands

The device commands provide low level access to the device functionality for a
single _UHPPOTE_ controller. Commands that require a device ID accept either the controller serial number or the controller name if defined defined in the _uhppoted.conf_ file e.g.
```
uhppote-cli get-device 405419896
uhppote-cli get-device alpha
```
Controller names are not case sensitive but should be enclosed in quotes if they contain spaces.

#### `get-devices`

Retrieves a list of the controllers accessible on the local LAN (i.e. can receive and respond to UDP broadcasts). The command returns a fixed width columnar table with:

- (`name`)
- `serial number`
- `IP address`
- `subnet mask`
- `gateway IP address`
- `MAC address`
- `firmware version`
- `current date`
```
uhppote-cli [options] get-devices

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Example:

  uhppote-cli get-devices
  
         201020304  192.168.1.101  255.255.255.0  192.168.1.1  52:fd:fc:07:21:82  v6.62  2020-01-01
  Beta   303986753  192.168.1.100  255.255.255.0  192.168.1.1  52:fd:fc:07:21:82  v8.92  2019-08-15
  Alpha  405419896  192.168.1.100  255.255.255.0  192.168.1.1  00:12:23:34:45:56  v8.92  2018-11-05
```
**NOTE**
1. The `name` field is retrieved from the _uhppoted.conf_ file. 
2. The `name` column is omitted entirely if none of the devices has a name defined in the _uhppoted.conf_ file.

#### `get-device`

Retrieves the controller information for a single controller accessible on the local LAN (i.e. can receive and respond to UDP broadcasts) or configured in the communal `uhppoted.conf` configuration (or custom configuration, if specified). The command returns a fixed width columnar table with:

- (`name`)
- `serial number`
- `IP address`
- `subnet mask`
- `gateway IP address`
- `MAC address`
- `firmware version`
- `current date`
```
uhppote-cli [options] get-devices

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Example:

  uhppote-cli --debug --conf ./405419896.conf --timeout 1.25s get-device 405419896
  
  ... request
  ...          00000000  17 94 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
  ...          00000010  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
  ...          00000020  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
  ...          00000030  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
  ...
  ... sent 64 bytes to 192.168.1.100:60000
  ... received 64 bytes from 192.168.1.100:60000
  ...          00000000  17 94 00 00 78 37 2a 18  c0 a8 01 64 ff ff ff 00  |....x7*....d....|
  ...          00000010  c0 a8 01 01 00 12 23 34  45 56 08 92 20 20 05 21  |......#4EV..  .!|
  ...          00000020  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
  ...          00000030  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
  ...

  Alpha  405419896  192.168.1.100   255.255.255.0   192.168.1.1     00:12:23:34:45:56 v8.92 2018-11-05
```
**NOTE**
1. The `name` field is retrieved from the _uhppoted.conf_ file. 
2. The `name` field is omitted entirely if the devices does not have a name defined in the _uhppoted.conf_ file.

#### `set-address`

Sets the controller IP address.

```
uhppote-cli [options] set-address <device> <address> [mask] [gateway]

  <device>      (required) Controller serial number
  <address>     (required) IP address assigned to controller
  <mask>        (optional) Subnet mask. Defaults to 255.255.255.0 if not provided.
  <gateway>     (optional) Gateway IP address. Defaults to 0.0.0.0 if not provided.

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the controllers

  Example:

  uhppote-cli set-address 405419896 192.168.1.49 255.255.255.0 192.168.1.1
```

#### `get-time`

Retrieves the current date/time for a single controller accessible on either the local LAN or configured in the communal `uhppoted.conf` configuration (or custom configuration, if specified). The command returns a fixed width columnar table with:

- `serial number`
- `date` _(YYYY-MM-DD)_
- `time` _(HH:mm:ss)_
```
uhppote-cli [options] get-time <device>

  <device>      (required) Controller serial number

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Example:

  uhppote-cli get-time 405419896
  
  405419896  2020-05-21 09:49:56
```

#### `set-time`

Sets the current date/time for a single controller accessible on either the local LAN or configured in the communal `uhppoted.conf` configuration (or custom configuration, if specified). The command returns a fixed width columnar table with:

- `serial number`
- `date` _(YYYY-MM-DD)_
- `time` _(HH:mm:ss)_
```
uhppote-cli [options] set-time <device> [date/time]

  <device>      (required) Controller serial number
  <date/time>   Date/time in YYYY-MM-DD HH:mm:ss format or use 'now' to use the current time on the uhppote-cli host). Defaults to 'now'.

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Examples:

  uhppote-cli set-time 405419896 now
  uhppote-cli set-time 405419896 2020-05-21 10:55:32
  
  405419896  2020-05-21 10:55:32
```

#### `get-door-delay`

Retrieves the current door delay setting (open duration) for a single controller accessible on either the local LAN or configured in the communal `uhppoted.conf` configuration (or custom configuration, if specified). The command returns a fixed width columnar table with:

- `serial number`
- `door` _[1..4]_
- `delay` _(seconds)_
```
uhppote-cli [options] get-door-delay <device> <door>

  <device>      (required) Controller serial number
  <door>        (required) Door number (1..4)

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Example:

  uhppote-cli get-door-delay 405419896 3
  
  405419896  3 7
```

#### `set-door-delay`

Sets the door delay setting (open duration) for a single controller accessible on either the local LAN or configured in the communal `uhppoted.conf` configuration (or custom configuration, if specified). The command returns a fixed width columnar table with:

- `serial number`
- `door` _[1..4]_
- `delay` _(seconds)_
```
uhppote-cli [options] get-door-delay <device> <door>

  <device>      (required) Controller serial number
  <door>        (required) Door number (1..4)
  <delay>       (required) Delay (in seconds)

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Example:

  uhppote-cli set-door-delay 405419896 3 5
  
  405419896  3 5
```

#### `get-door-control`

Retrieves the current door control setting for a single controller accessible on either the local LAN or configured in the communal `uhppoted.conf` configuration (or custom configuration, if specified). The command returns a fixed width columnar table with:

- `serial number`
- `door` _[1..4]_
- `control code` _[1..3]_
- `control state` _(normally open/normally closed/controlled)_
```
uhppote-cli [options] get-door-control <device> <door>

  <device>      (required) Controller serial number
  <door>        (required) Door number (1..4)

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Example:

  uhppote-cli get-door-control 405419896 3
  
  405419896  3 3 (controlled)
```

#### `set-door-control`

sets the current door control setting for a single controller accessible on either the local LAN or configured in the communal `uhppoted.conf` configuration (or custom configuration, if specified). The command returns a fixed width columnar table with:

- `serial number`
- `door` _[1..4]_
- `control code` _[1..3]_
- `control state` _(normally open/normally closed/controlled)_
```
uhppote-cli [options] set-door-control <device> <door> <state>

  <device>      (required) Controller serial number
  <door>        (required) Door number (1..4)
  <control>     (required) normally open/normally closed/controlled

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Example:

  uhppote-cli set-door-control 405419896 3 "normally open"
  
  405419896  3 1 (normally open)
```

#### `get-listener`

Retrieves the current IP address + port assigned to receive event notifications from a controller. The command returns a fixed width columnar table with:

- `serial number`
- `listener ` _(IP address:port)_
```
uhppote-cli [options] get-listener <device>

  <device>      (required) Controller serial number

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Example:

  uhppote-cli get-listener 405419896
  
  405419896  192.168.1.100:60001
```

#### `set-listener`

Sets the IP address + port to which to send controller event notifications. The command returns a fixed width columnar table with:

- `serial number`
- `sucess` _(true/false)_
```
uhppote-cli [options] set-listener <device> <address:port>

  <device>       (required) Controller serial number
  <address:port> (required) Listener IP address and port

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Example:

  uhppote-cli set-listener 405419896 192.168.1.100:60001
 
  405419896  true
```

#### `record-special-events`

Enables or disables events for door open, door closed and door button pressed for a single controller accessible on either the
local LAN or configured in the communal `uhppoted.conf` configuration (or custom configuration, if specified). The command 
returns the result as a fixed width columnar table with:

- `serial number`
- `succeeded` _true/false_

```
uhppote-cli [options] record-special-events <device> <enable>

  <device>      (required) Controller serial number
  <enable>      (optional) Enables or disables door open, closed and button pressed events. Defaults to `true`

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Example:
  > uhppote-cli record-special-events 405419896 false
    405419896  true
```

#### `get-status`

Retrieves the controller status for a single controller accessible on the local LAN (i.e. can receive and respond to UDP broadcasts) or configured in the communal `uhppoted.conf` configuration (or custom configuration, if specified). The command returns a fixed width columnar table with:

- `device ID`
- `door states` [1..4]
- `door button states` [1..4]
- `system state`
- `system date-time`
- `packet sequence number`
- `special info`
- `relay state`
- `input state`
- `last event: index`
- `last event: type`
- `last event: access granted`
- `last event: door`
- `last event: door opened`
- `last event: user ID`
- `last event: timestamp`
- `last event: result code`

```
uhppote-cli [options] get-status <device>

  <device>      (required) Controller serial number

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Example:

  uhppote-cli get-status 405419896
  
  405419896  false false false false false false false false 0    2021-04-24 09:11:17 0          0 00 00 | 69    2   true  1 1     0          2019-08-10 10:28:32 44
```

**NOTE**
The event fields are separated from the static data by a '|' and are not displayed if the controller does not have a valid 'last event'.

#### `get-cards`

Retrieves all access card records from a controller. A card record comprises:

- `card number`
- `start date` _date from which the card is active_
- `end date`   _date after which the card is inactive_
- `door1`     _Y,N or associated time profile for door 1_
- `door2`     _Y,N or associated time profile for door 2_
- `door3`     _Y,N or associated time profile for door 3_
- `door4`     _Y,N or associated time profile for door 4_
```
uhppote-cli [options] get-cards <device ID>

  <device ID>   (required) Controller serial number (or name)

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Example:

  uhppote-cli get-cards 405419896
  
  8165537  2021-01-01 2021-12-31 Y N N N
  8165539  2021-01-01 2021-12-31 N N N N
  8165538  2021-01-01 2021-12-31 Y N Y 29
```

#### `get-card`

Retrieves a single access card record from a controller. A card record comprises:

- `card number`
- `start date` _date from which the card is active_
- `end date`   _date after which the card is inactive_
- `door1`     _Y,N or associated time profile for door 1_
- `door2`     _Y,N or associated time profile for door 2_
- `door3`     _Y,N or associated time profile for door 3_
- `door4`     _Y,N or associated time profile for door 4_
```
uhppote-cli [options] get-card <device ID> <card number>

  <device ID>   (required) Controller serial number (or name)
  <card number> (required) Card number of card to be retrieved

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Example:

  uhppote-cli get-card 405419896 8165538
  
  8165538  2021-01-01 2021-12-31 Y N Y 29
```

#### `put-card`

Creates (or updates) an access card record on a controller, with the following information:

- `card number`
- `start date` _date from which the card is active_
- `end date`   _date after which the card is inactive_
- `door1`     _Y,N or associated time profile for door 1_
- `door2`     _Y,N or associated time profile for door 2_
- `door3`     _Y,N or associated time profile for door 3_
- `door4`     _Y,N or associated time profile for door 4_
```
uhppote-cli [options] put-card <device ID> <card number> <start> <end> <doors>

  <device ID>   (required) Controller serial number (or name)
  <card number> (required) Access card number
  <start>       (required) Start date from which the card is enabled, formatted as YYYY-mm-dd
  <end>         (required) End dates after which the card is no longer enabled, formatted as YYYY-mm-dd
  <doors>       (optional) Comma separated list of doors for which the card grants access. Time profiled access for a door can be specified as door:profile.

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Example:

  uhppote-cli put-card 405419896 8165538 2021-01-01 2021-12-31 1,3,4:29
  405419896 8165538 true
```
**NOTES**
1. For a door associated with a time profile, `uhppote-cli` requires the time profile to be an existing time profile defined in the controller. The controller itself does not enforce this requirement but linking to a non-existent time profile counter-intuitively seems to allow access at any time of day.

#### `delete-card`
Unconditionally deletes an access card record from a controller, returning `true` if successful.
```
uhppote-cli [options] delete-card <device ID> <card number>

  <device ID>   (required) Controller serial number (or name)
  <card number> (required) Card number of card to be deleted

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Example:

  uhppote-cli delete-card 405419896 8165538
  
  405419896 8165538 true
```

#### `delete-all`

Unconditionally deletes all cards from a controller, returning `true` if successful.
```
uhppote-cli [options] delete-all <device ID>

  <device ID>   (required) Controller serial number (or name)

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Example:

  uhppote-cli delete-all 405419896
  405419896 true
```

#### `get-time-profile`

Retrieves a time profile from a controller, returning a (space delimited) time profile:

- `serial number`
- `profile ID`
- `from:to`
- `weekdays`
- `time segments 1,2 and 3`
- `linked profile ID`
```
uhppote-cli [options] get-time-profile <device ID> <profile ID>

  <device ID>   (required) Controller serial number (or name)
  <profile ID>  (required) Time profile ID (in the interval [2..254])

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Example:

  uhppote-cli get-time-profile 405419896 2
  405419896  2 2021-04-01:2021-12-31 Mon,Wed,Fri 08:30-11:30,13:45-17:00 3
  
  uhppote-cli get-time-profile 405419896 29
  405419896 29 NO ACTIVE TIME PROFILE
```

#### `set-time-profile`

Creates (or updates) a controller time profile. A time profile comprises the following information:

- `serial number`
- `profile ID`
- `from:to`
- `weekdays`
- `time segments 1-3`
- `linked profile ID`
```
uhppote-cli [options] get-time-profile <device ID> <profile ID> <from:to> <weekdays> <segments> <linked>

  <device ID>   (required) Controller serial number (or name)
  <profile ID>  (required) Time profile ID (in the interval [2..254])
  <from:to>     (required) Start and end dates for the time profile ID, formatted as YYYY-mm-dd:YYYY-mm-dd
  <weekdays>    (optional) Comma separated list of weekdays for which the time profile is active. If omitted, the time profile is not active on any weekday.
  <segments>    (optional) Comma separated list of up to 3 time segments formatted as HH:mm-HH:mm. Omitted segments are created as 00:00-00:00.
  <linked>      (optional) Time profile ID to be linked to this time profile to allow more than 3 segments to be associated with a time profile. 

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Example:

  uhppote-cli set-time-profile 405419896 29 2021-04-01:2021-12-31 3 Mon,Wed,Fri 08:30-11:30,,13:45-17:00 3
  405419896: time profile 29 created
```
**NOTES**
1. `uhppote-cli` requires the linked profile to be an existing time profile defined in the controller. The controller itself does not enforce this requirement but linking to a non-existent time profile counter-intuitively seems to allow access at any time of day.
2. `uhppote-cli` does not allow a time profile to define a link that creates a circular chain of time profiles. Again, the controller itself does not enforce this requirement but is indicative of a mistake in the definition of the time profiles.

#### `clear-time-profiles`

Deletes all time profiles from a controller.

```
uhppote-cli [options] clear-time-profiles <device ID>

  <device ID>   (required) Controller serial number (or name)

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Example:

  uhppote-cli clear-time-profiles 405419896
  405419896 true
```
**NOTES**
1. Clearing all time profiles will allow access at any time of day to cards that were previously _time managed_. In general this command is intended to be used when redefining all the time profiles.

#### `get-time-profiles`

Retrieves all time profiles from a controller, optionally storing them in a TSV file.

```
uhppote-cli [options] get-time-profiles <device ID> <TSV>

  <device ID>   (required) Controller serial number (or name)
  <TSV>         (optional) TSV file to which to write retrieved profiles.

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Example:

  uhppote-cli get-time-profiles 405419896
  
  -------------------------------------------
  TIME PROFILES 405419896 2021-05-13 10:36:21
  -------------------------------------------
  Profile  From       To          Mon Tue Wed Thurs Fri Sat Sun  Start1 End1   Start2 End2   Start3 End3   Linked
  2        2021-04-01 2021-10-29  N   Y   N   N     Y   N   Y    08:30  11:30  00:00  03:00  13:45  17:00  3  
  3        2021-04-02 2021-11-30  N   N   Y   N     Y   Y   Y    09:31  11:31  01:00  04:00  13:46  17:01     
  29       2021-04-03 2021-12-31  N   N   N   Y     N   Y   Y    10:32  11:32  02:00  05:00  13:47  17:02     
  75       2021-04-01 2021-12-31  Y   N   Y   N     Y   N   N    08:30  11:30  00:00  00:00  13:45  17:00     
  
  uhppote-cli get-time-profiles 405419896 405419896.tsv
```

#### `set-time-profiles`

Creates (or updates) time profiles on a controller from a TSV file. 

Invalid time profiles (e.g. with missing or otherwise invalid _from_ and _to_ dates or invalid segments) are ignored with a warning. Likewise, time profiles that link to an undefined time profile or time profiles with a linked profile that would create a circular chain are ignored with a warning. 
```
uhppote-cli [options] set-time-profiles <device ID> <TSV>

  <device ID>   (required) Controller serial number (or name)
  <TSV>         (required) TSV file with the time profiles to be created (or updated)

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Sample TSV file format:

Profile From        To          Mon Tue Wed Thurs Fri Sat Sun Start1  End1   Start2  End2  Start3  End3  Linked
2       2021-04-01  2021-10-29  N   Y   N    N    Y   N   Y   08:30   11:30  00:00   03:00 13:45   17:00 3
3       2021-04-02  2021-11-30  N   N   Y    N    Y   Y   Y   09:31   11:31  01:00   04:00 13:46   17:01 
29      2021-04-03  2021-12-31  N   N   N    Y    N   Y   Y   10:32   11:32  02:00   05:00 13:47   17:02 
75      2021-04-01  2021-12-31  Y   N   Y    N    Y   N   N   08:30   11:30  00:00   00:00 13:45   17:00 
100     2021-04-01  2021-12-31  Y   N   Y    N    Y   N   N   08:30   11:30  00:00   00:00 13:45   17:00 
101     2021-04-01  2021-10-29  N   Y   N    N    Y   N   Y   08:30   11:30  00:00   03:00 13:45   17:00 102

  Example:

  uhppote-cli set-time-profiles 405419896 405419896.tsv
  
  ./bin/uhppote-cli set-time-profiles 405419896 405419896.tsv
   ... set time profile 3
   ... set time profile 29
   ... set time profile 105
   ... set time profile 2

   WARN  profile 51  - linked time profile 50 is not defined
   WARN  profile 75  - linking to time profile 75 creates a circular reference [75 75]
   WARN  profile 150 - 'To' date (2021-04-01) is before 'From' date (2021-12-31)
   WARN  profile 151 - segment 1 'End' (09:31) is before 'Start' (11:31)  
```
**NOTES** 
1. `set-time-profiles` does not clear existing time profiles from the controller i.e. although not recommended, profiles in the file can link to individually defined existing profiles.
2. There is no requirement for the profiles to be in any particular order e.g. `uhppote-cli` is capable of creating a time profile that is linked to a profile that is defined after it in the TSV file.


#### `clear-task-list`

Deletes all defined tasks from a controller.

```
uhppote-cli [options] clear-task-list <device ID>

  <device ID>   (required) Controller serial number (or name)

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Example:

  uhppote-cli clear-task-list 405419896
  405419896 true
```

#### `refresh-task-list`

Activates all tasks created by `add-task`.

```
uhppote-cli [options] refresh-task-list <device ID>

  <device ID>   (required) Controller serial number (or name)

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Example:

  uhppote-cli refresh-task-list 405419896
  405419896 true
```

#### `add-task`

Create a new task on the controller. The task will only be activated after invoking `refresh-task-list`. A task definition comprises:

- `serial number`
- `task ID`
- `door`
- `from:to`
- `weekdays`
- `start time`
- `number of cards allowed for 'more cards'`

```
uhppote-cli [options] add-task <device ID> <task> <door> <from:to> <days> <start> <cards>

  <device ID>   (required) Controller serial number (or name)
  <task>        (required) Task ID or name, corresponding to one of the predefined operations
                           listed below
  <door>        (required) Door ID [1..4] to which the operation applies.
  <from:to>     (required) Start and end dates between which the task is active, formatted as YYYY-mm-dd:YYYY-mm-dd
  <weekdays>    (optional) Comma separated list of weekdays for which the task is active. Defaults to 'all' if omitted.
  <start>       (optional) Time from which task is active, formatted as HH:mm. Defaults to 00:00 if omitted.
  <cards>       (optional) Number of 'more cards' for the 'enable more cards' operation.
   
  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Tasks: 
    1   control door
    2   unlock door
    3   lock door
    4   disable time profile
    5   enable time profile
    6   enable card, no password
    7   enable card+IN password
    8   enable card+password
    9   enable more cards
    10  disable more cards
    11  trigger once
    12  disable pushbutton
    13  enable pushbutton
  
  Example:

  uhppote-cli add-task 405419896 3 4 2021-01-01:2021-12-31 Mon,Fri 08:30
  405419896 true

  uhppote-cli add-task 405419896 'enable more cards' 4 2021-01-01:2021-12-31 Mon,Fri 08:30 29
  405419896 true
```

#### `set-task-list`

Sets the task list on a controller from a TSV file. The command clears the active task list, adds all tasks
defined in the TSV file and then invokes `refresh-task-list` to activate the added tasks.

```
uhppote-cli [options] set-task-list <device ID> <file>

  <device ID>   (required) Controller serial number (or name)
  <file>        (required) TSV containing the task list for the controller
   
  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Sample TSV file format:
  
   Task               Door  From        To          Mon Tue Wed Thurs Fri Sat Sun Start Cards
   1                  3     2021-04-01  2021-10-29  N   Y   N    N    Y   N   Y   08:30 0
   disable pushbutton 3     2021-04-02  2021-11-30  N   N   Y    N    Y   Y   Y   09:45 0
   8                  3     2021-04-03  2021-12-31  N   N   N    Y    N   Y   Y   10:15 29
  
  Example:

  uhppote-cli set-task-list 405419896 405419896.tasks

  ./bin/uhppote-cli set-task-list 405419896 405419896.tasks
   ... 405419896 cleared task list
   ... created task defintion 1
   ... created task defintion 2
   ... created task defintion 3
   ... 405419896 refreshed task list
```


#### `get-events`

Retrieves the start and end range of the events stored on a controller as well as the 
current event index.

```
uhppote-cli [options] get-events <device ID>

  <device ID>   (required) Controller serial number (or name)

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Examples:
  > uhppote-cli get-events 405419896
    405419896  1  69 53
  
  > uhppote-cli get-events 303986753
    303986753  NO EVENTS
```

#### `get-event`

Retrieves the record for a single event from a controller, comprising:
- `device ID`
- `event ID`
- `timestamp`
- `card number`
- `door`
- `access granted`
- `reason code`

The event ID should be within the event range returned by `get-events`.
```
uhppote-cli [options] get-event <device ID> <event ID>

  <device ID>   (required) Controller serial number (or name)
  <event ID>    (optional) ID of event to be retrieved. If omitted, the event at the current event index is returned 
                           and the event index is incremented. `first`, `last`, `current` and `next` retrieve the 
                           _first_, _last_, _current_ and _next_ stored events respectively. The controller _current event
                           index_ is only incremented for `next`.

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Examples:
  > uhppote-cli get-event 405419896
    405419896  25     2019-07-24 20:12:40 98765432 3 true  0 
    
  > uhppote-cli get-event 405419896 27
    405419896  27     2019-07-24 20:12:47 98765432 3 true  0     
    
  > uhppote-cli get-event 405419896 first
    405419896  23      2019-07-24 20:12:33 98765432 3 true  0     
    
  > uhppote-cli get-event 405419896 last
    405419896  67     2019-07-24 20:12:43 98765432 4 false 6     
    
  > uhppote-cli get-event 405419896 current
    405419896  24     2019-07-24 20:12:43 98765432 4 false 6     
    
  > uhppote-cli get-event 405419896 next
    405419896  25     2019-07-24 20:12:43 98765432 4 false 6     

  > uhppote-cli get-event 405419896 17263
    ERROR: 405419896:  no event at index 17263

  > uhppote-cli get-event 405419896 17263 17
    ERROR: 405419896:  event at index has been overwritten
```

#### `get-event-index`

Retrieves the current event index from a controller. The event index is a convenient user value that is typically used to store the ID of the most recently retrieved event.

```
uhppote-cli [options] get-event-index <device ID>

  <device ID>   (required) Controller serial number (or name)

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Examples:
  > uhppote-cli get-event-index 405419896
    405419896  24
```

#### `set-event-index`

Sets the current event index on a controller. The event index is a convenient user value that is typically used to store the ID of the most recently retrieved event. It is not automatically managed by the controller.

```
uhppote-cli [options] set-event-index <device ID> <index>

  <device ID>   (required) Controller serial number (or name)
  <index>       (required) ID of event to set as the event index.

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --debug       Displays verbose debugging information, in particular the communications with the  controllers

  Examples:
  > uhppote-cli set-event-index 405419896 19
    405419896  19       true
```

#### `open`
Unconditionally unlocks a door, provided the door control state is _controlled_ i.e. not _normally open_ or _normally closed_.
```
uhppote-cli [options] open <device ID> <door>

  <device ID>   (required) Controller serial number (or name)
  <door>        (required) ID of door to unlock ([1..4])

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Examples:
  > uhppote-cli open 405419896 3
    405419896  true
```

#### `listen`
Establishes a _listening_ socket on the _listen_ `address:port` for events sent from controllers. Each event record comprises the same information as that returned by `get-status`:

- `device ID`
- `door states` [1..4]
- `door button states` [1..4]
- `system state`
- `system date-time`
- `packet sequence number`
- `special info`
- `relay state`
- `input state`
- `last event: index`
- `last event: type`
- `last event: access granted`
- `last event: door`
- `last event: door opened`
- `last event: user ID`
- `last event: timestamp`
- `last event: result code`

```
uhppote-cli [options] listen

  <device ID>   (required) Controller serial number (or name)

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Examples:
  > uhppote-cli listen
    405419896  false false false false false false false false 0    2021-05-14 11:14:18 0          0 00 00 | 71    1   false 3 1     8165538    2021-05-14 11:14:18 6
    405419896  false false false false false false false false 0    2021-05-14 11:14:21 0          0 00 00 | 72    1   false 3 1     8165538    2021-05-14 11:14:21 6
    405419896  false false false false false false false false 0    2021-05-14 11:14:24 0          0 00 00 | 73    1   false 3 1     8165538    2021-05-14 11:14:24 6

```

### ACL commands

The ACL (_access control list_) commands manage access permissions across the set of _UHPPOTE_ controllers configured in the `conf` file. The following commands are supported:

- `grant`
- `revoke`
- `show`
- `load-acl`
- `get-acl`
- `compare-acl`

### ACL file format

The only currently supported ACL file format is TSV (tab separated values) and is expected to be formatted as follows:

    Card Number From  To  Workshop  Side Door Front Door  Garage  Upstairs  Downstairs  Tower Cellar
    123465537 2020-01-01  2020-12-31  N N Y N Y N Y Y
    231465538 2020-01-01  2020-12-31  Y N Y N N 29 N N
    635465539 2020-01-01  2020-12-31  N N N N Y N Y Y

| Field         | Description                                                                |
|---------------|----------------------------------------------------------------------------|
| `Card Number` | Access card number                                                         |
| `From`        | Date from which card is valid (_valid from 00:00 on that date_)            |
| `To`          | Date until which card is valid (_valid until 23:59 on that date_)          |
| `<door>`      | Door name matching controller configuration (_case and space insensitive_) |
| `<door>`      | ...                                                                        |
| ...           |                                                                            |

The ACL file must include a column for each controller + door configured in the _devices_ section of the `uhppoted.conf` file used to configure the utility. Permissions for a door can be:
- Y
- N
- time profile ID (e.g. 29)

An [example ACL file](https://github.com/uhppoted/uhppoted/blob/master/runtime/simulation/405419896.acl) is included in the full `uhppoted` distribution, along with the matching [_conf_](https://github.com/uhppoted/uhppoted/blob/master/runtime/simulation/405419896.conf) file.

#### `grant`

Grants access permissions to a single card across the set of configured UHPPOTE controllers. The `grant` command extends
the functionality of the device `put-card` command to set the access permissions across multiple controllers using door 
names rather than the device specific controller + door ID's. The _granted_ permissions are added to the existing card
access permissions.

*NOTE: The date range during which a card has access is automatically widened to accommodate the earliest _from_ date and latest
_to_ date across the controllers.*

```
uhppote-cli [options] grant <card> <from> <to> [doors]

  <card>        Card number to be granted access
  <from>        Date from which the card is granted access, in yyyy-mm-dd format. Access is 
                granted from 00:00 on the 'from' date.
  <to>          Date until which the card is granted access, , in yyyy-mm-dd format. Access
                is granted until 23:59 on the 'to' date.

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Example:

  uhppote-cli grant 918273645 2020-01-01 2020-12-31 Front Door, Workshop

```

#### `revoke`

Revokes access permissions for a single card across the set of configured UHPPOTE controllers. The `revoke` command 
extends the functionality of the device `put-card` command to revoke the access permissions across multiple controllers 
using door names rather than the device specific controller + door ID's. The _revoked_ permissions are removed from the 
existing card access permissions.

```
  uhppote-cli [options] revoke <card> [doors]

  <card>        Card number to be granted access

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Example:

  uhppote-cli revoke 918273645 Garage

```

#### `show`

Retrieves and displays the access permissions for a single card across the set of configured UHPPOTE controllers. 

```
   uhppote-cli [options] show <card>

  <card>        Card number for which access permissions should be retrieved

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Example:

  uhppote-cli show 918273645

```

#### `load-acl`

Loads the access permissions from an ACL file to the set of configured UHPPOTE controllers. A sample [ACL](https://github.com/uhppoted/uhppoted/blob/master/runtime/simulation/405419896.acl) file is included in the full 
`uhppoted` distribution.

```
   uhppote-cli [options] load-acl <ACL file>

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Example:

  uhppote-cli --debug --conf warehouse.conf load-acl warehouse.acl

```

### `get-acl`

Fetches the cards stored in the set of configured UHPPOTE controllers, creates a matching ACL file from the UHPPOTED controller configuration and writes it to a TSV file. 

```
  uhppote-cli [options] get-acl <file>

  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Example:

  uhppote-cli --debug --conf warehouse.conf get-acl 2020-05-18.acl
```

### `compare-acl`

Compares the cards stored on the set of configured UHPPOTE controllers with an authoritative ACL file and generates an
exception report.

```
   uhppote-cli [options] compare-acl <file> <report>

   <report>     Optional output file for the exception report. The report is printed
                to the console if a report file is not supplied.
  
  Options: 
  --config      Sets the uhppoted.conf file to use for controller configurations
  --bind        Overrides the default (or configured) bind IP address for a command
  --broadcast   Overrides the default (or configured) broadcast IP address to which to send a command
  --listen      Overrides the default (or configured) listen IP address on which to listen for events
  --timeout     Sets the timeout for a response from a controller (default value is 2.5s)
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Example:

  uhppote-cli --debug --conf warehouse.conf compare-acl warehouse.acl 2020-05-18.rpt
```


