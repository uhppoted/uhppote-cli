# uhppote-cli

CLI for the *UHPPOTE UT0311-L0x* TCP/IP Wiegand access control boards. 

Supported operating systems:
- Linux
- MacOS
- Windows
- ARM7

## Raison d'être

The manufacturer supplied application is _Windows-only_ and provides limited support for integration with other
systems.

## Releases

| *Version* | *Description*                                                                                            |
| --------- | -------------------------------------------------------------------------------------------------------- |
| v0.6.3    | Reworked get-cards to handle deleted records                                                             |
| v0.6.2    | Fixed get-events for controllers without any retrievable events and improved configuration file handling |
| v0.6.1    | Added ACL commands to simplify managing card permissions across multiple controllers                     |
| v0.6.0    | Maintenance release to keep compatibility with updated `uhppote-core`                                    |
| v0.5.1    | Initial release following restructuring into standalone Go *modules* and *git submodules*                |

## Installation

Executables for all the supported operating systems are packaged in the [releases](https://github.com/uhppoted/uhppote-cli/releases). The provided archives contain the executables for all the operating systems - operating system specific tarballs can be found in the [uhppoted](https://github.com/uhppoted/uhppoted/releases) releases.

Installation is straightforward - download the archive and extract it to a directory of your choice and place the executable in a directory in your PATH. 

### `uhppoted.conf`

By default, `uhppote-cli` uses the communal `uhppoted.conf` configuration file shared by all the `uhppoted` project modules which is (or will eventually be) documented in [uhppoted](https://github.com/uhppoted/uhppoted). 

A valid configuration file is required **only** for system configurations where controllers are not findable on the local LAN (i.e. cannot receive and reply to UDP broadcasts) or for use with the _ACL_ commands which require a valid `uhppoted.conf` to resolve the door ID to controller + door number. For _device_ commands the configuration file will used if present otherwise the internal default configuration will be used.

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
go build -o bin ./...
```

The above commands build the `uhppote-cli` executable to the `bin` directory.

#### Dependencies

| *Dependency*                                                                 | *Description*                              |
| ---------------------------------------------------------------------------- | ------------------------------------------ |
| [com.github/uhppoted/uhppote-core](https://github.com/uhppoted/uhppote-core) | Device level API implementation            |
| [com.github/uhppoted/uhppoted-api](https://github.com/uhppoted/uhppoted-api) | common API for external applications       |
| golang.org/x/lint/golint                                                     | Additional *lint* check for release builds |

## uhppote-cli

Usage: ```uhppote-cli [options] <command> <parameters>```

General_ commands:

- `help`
- `version`

Device commands:

- `get-devices`
- `get-device`
- `set-address`
- `get-status`
- `get-time`
- `set-time`
- `get-door-delay`
- `set-door-delay`
- `get-door-control`
- `set-door-control`
- `get-listener`
- `set-listener`
- `get-cards`
- `get-card`
- `put-card`
- `delete-card`
- `delete-all`
- `get-events`
- `get-swipe-index`
- `set-event-index`
- `open`
- `listen`

ACL commands:

- `grant`
- `revoke`
- `load-acl`
- `get-acl`
- `compare-acl`

Common command options:
```
  --config      Sets the uhppoted.conf file to use for controller configurations
  --debug       Displays verbose debugging information, in particular the 
                communications with the UHPPOTE controllers

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
single _UHPPOTE_ controller.

#### `get-devices`

Retrieves a list of the controllers accessible on the local LAN (i.e. can receive and respond to UDP broadcasts). The command returns a fixed width columnar table with:

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
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Example:

  uhppote-cli get-devices
  
  999        192.168.1.100   255.255.255.0   192.168.1.1     00:12:23:34:45:56 0892 2020-05-21
  303986753  192.168.1.100   255.255.255.0   192.168.1.1     52:fd:fc:07:21:82 0892 2020-05-21
  405419896  192.168.1.100   255.255.255.0   192.168.1.1     00:12:23:34:45:56 0892 2020-05-21

```

#### `get-device`

Retrieves the controller information for a single controller accessible on the local LAN (i.e. can receive and respond to UDP broadcasts) or configured in the communal `uhppoted.conf` configuration (or custom configuration, if specified). The command returns a fixed width columnar table with:

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
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers

  Example:

  uhppote-cli --debug --conf ./405419896.conf  get-device 405419896
  
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

  405419896  192.168.1.100   255.255.255.0   192.168.1.1     00:12:23:34:45:56 0892 2020-05-21
```

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
  --debug       Displays verbose debugging information, in particular the communications with the controllers

  Example:

  uhppote-cli set-address 405419896 192.168.1.49 255.255.255.0 192.168.1.1
```

#### `get-status`

Retrieves the controller status for a single controller accessible on the local LAN (i.e. can receive and respond to UDP broadcasts) or configured in the communal `uhppoted.conf` configuration (or custom configuration, if specified). The command returns a fixed width columnar table with:

- `serial number`
- `last event: index`
- `last event: type`
- `last event: access granted`
- `last event: door`
- `last event: door opened`
- `last event: user ID`
- `last event: timestamp`
- `last event: result code`
- `door states` [1..4]
- `door button states` [1..4]
- `system state`
- `system date-time`
- `packet number` _(? TBC)_
- `backup status` _(? TBC)_
- `special message` _(? TBC)_
- `battery status` _(? TBC)_
- `fire alarm status` _(? TBC)_
```
uhppote-cli [options] get-status <device>

  <device>      (required) Controller serial number

  Options: 

  --config      Sets the uhppoted.conf file to use for controller configurations
  --debug       Displays verbose debugging information, in particular the communications with the  controllers

  Example:

  uhppote-cli get-status 405419896
  
  405419896  69    2   true  1 true  3922570474 2019-08-10 10:28:32 44  false false false false false false false false 0 2020-05-21 09:42:21 0 0 0 0 0
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
  --debug       Displays verbose debugging information, in particular the communications with the controllers

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
  --debug       Displays verbose debugging information, in particular the communications with the controllers

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
  --debug       Displays verbose debugging information, in particular the communications with the controllers

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
  --debug       Displays verbose debugging information, in particular the communications with the controllers

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
  --debug       Displays verbose debugging information, in particular the communications with the controllers

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
  --debug       Displays verbose debugging information, in particular the communications with the controllers

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
  --debug       Displays verbose debugging information, in particular the communications with the controllers

  Example:

  uhppote-cli get-listener 405419896
  
  405419896  192.168.1.100:60001
```

#### `set-listener`

Sets the IP address + port to which to send controller event notifications. The command returns a fixed width columnar table with:

- `serial number`
- `sucsess` _(true/false)_
```
uhppote-cli [options] set-listener <device> <address:port>

  <device>       (required) Controller serial number
  <address:port> (required) Listener IP address and port

  Options: 

  --config      Sets the uhppoted.conf file to use for controller configurations
  --debug       Displays verbose debugging information, in particular the communications with the controllers

  Example:

  uhppote-cli set-listener 405419896 192.168.1.100:60001
 
  405419896  true
```
#### `get-cards`

#### `get-card`

#### `put-card`

#### `delete-card`

#### `delete-all`

#### `get-events`

#### `get-swipe-index`

#### `set-event-index`

#### `open`

#### `listen`

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
    231465538 2020-01-01  2020-12-31  Y N Y N N Y N N
    635465539 2020-01-01  2020-12-31  N N N N Y N Y Y

| Field         | Description                                                                |
|---------------|----------------------------------------------------------------------------|
| `Card Number` | Access card number                                                         |
| `From`        | Date from which card is valid (_valid from 00:00 on that date_)            |
| `To`          | Date until which card is valid (_valid until 23:59 on that date_)          |
| `<door>`      | Door name matching controller configuration (_case and space insensitive_) |
| `<door>`      | ...                                                                        |
| ...           |                                                                            |

The ACL file must include a column for each controller + door configured in the _devices_ section of the `uhppoted.conf` file used to configure the utility.

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
  --debug       Displays verbose debugging information, in particular the communications 
                with the UHPPOTE controllers

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
  --debug       Displays verbose debugging information, in particular the communications 
                with the UHPPOTE controllers

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
  --debug       Displays verbose debugging information, in particular the communications 
                with the UHPPOTE controllers

  Example:

  uhppote-cli --debug --conf warehouse.conf load-acl warehouse.acl

```

### `store-acl`

Fetches the cards stored in the set of configured UHPPOTE controllers, and creates a matching ACL file from the UHPPOTED controller configuration and write it to a TSV file. 

```
  uhppote-cli [options] get-acl <file>

  Options: 

  --config      Sets the uhppoted.conf file to use for controller configurations
  --debug       Displays verbose debugging information, in particular the communications 
                with the UHPPOTE controllers

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
  --debug       Displays verbose debugging information, in particular the communications 
                with the UHPPOTE controllers

  Example:

  uhppote-cli --debug --conf warehouse.conf compare-acl warehouse.acl 2020-05-18.rpt
```


