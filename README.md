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

| *Version* | *Description*                                                                             |
| --------- | ----------------------------------------------------------------------------------------- |
| v0.6.1    | Added ACL commands to simplify managing card permissions across multiple controllers      |
| v0.6.0    | Maintenance release to keep compatibility with updated `uhppote-core`                     |
| v0.5.1    | Initial release following restructuring into standalone Go *modules* and *git submodules* |

## Installation

Executables for all the supported operating systems are packaged in the [releases](https://github.com/uhppoted/uhppote-cli/releases):

- [tar.gz](https://github.com/uhppoted/uhppote-cli/releases/download/v0.6.0/uhppote-cli_v0.6.1.tar.gz)
- [zip](https://github.com/uhppoted/uhppote-cli/releases/download/v0.6.0/uhppote-cli_v0.6.1.zip)

The above archives contain the executables for all the operating systems - operating system specific tarballs can be found in the
[uhppoted](https://github.com/uhppoted/uhppoted/releases) releases.

Installation is straightforward - download the archive and extract it to a directory of your choice and place the executable in a directory in your PATH. The `uhppote-cli` utility requires the following additional files:

- `uhppoted.conf`

### `uhppoted.conf`

`uhppoted.conf` is the communal configuration file shared by all the `uhppoted` project modules and is (or will 
eventually be) documented in [uhppoted](https://github.com/uhppoted/uhppoted). `uhppote-cli` requires the _devices_ 
section to resolve non-local controller IP addresses as well as door to controller door identities.

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

Supported commands:

- `help`
- `version`
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
```

### General

#### `help`

#### `version`

### Device commands

The device commands provide low level access to the device functionality for a
single _UHPPOTE_ controller.

#### `get-devices`

#### `get-device`

#### `set-address`

#### `get-status`

#### `get-time`

#### `set-time`

#### `get-door-delay`

#### `set-door-delay`

#### `get-door-control`

#### `set-door-control`

#### `get-listener`

#### `set-listener`

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

    Card Number	From	To	Workshop	Side Door	Front Door	Garage	Upstairs	Downstairs	Tower	Cellar
    123465537	2020-01-01	2020-12-31	N	N	Y	N	Y	N	Y	Y
    231465538	2020-01-01	2020-12-31	Y	N	Y	N	N	Y	N	N
    635465539	2020-01-01	2020-12-31	N	N	N	N	Y	N	Y	Y

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

Fetches the cards stored in the set of configured UHPPOTE controllers, and creates a matching ACL file from the UHPPOTED 
controller configuration and write it to a TSV file. 

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


