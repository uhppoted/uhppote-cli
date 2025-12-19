# CHANGELOG

## Unreleased

### Updated
1. Updated to Go 1.25.


## [0.8.11](https://github.com/uhppoted/uhppote-cli/releases/tag/v0.8.11) - 2025-07-01

### Added
1. `get/set-antipassback` command to manage the controller anti-passback mode.

### Updated
1. Updated to Go 1.24.


## [0.8.10](https://github.com/uhppoted/uhppote-cli/releases/tag/v0.8.10) - 2025-01-30

### Added
1. ARMv6 build target (RaspberryPi ZeroW).

### Updated
1. Added (optional) auto-send interval to _set-listener_ command.
2. Included auto-send interval (if not zero) in response from _get-listener_ command.


## [0.8.9](https://github.com/uhppoted/uhppote-cli/releases/tag/v0.8.9) - 2024-09-06

### Added
1. github nightly build with executable artifacts.
2. TCP/IP support.

### Updated
1. Changed default controller timezone from UTC to Local.
2. Updated to Go 1.23.


## [0.8.8](https://github.com/uhppoted/uhppote-cli/releases/tag/v0.8.8) - 2024-03-27

### Added
1. Implemented `restore-default-parameters` command to reset a controller to the manufacturer
   default configuration.

### Updated
1. Bumped Go version to 1.22.


## [0.8.7](https://github.com/uhppoted/uhppote-cli/releases/tag/v0.8.7) - 2023-12-01

### Added
1. Implemented `set-super-passwords` command to set the _super_ passwords for a door.

### Updated
1. Replaced `nil` event pointer with zero value in `get-status`.


## [0.8.6](https://github.com/uhppoted/uhppote-cli/releases/tag/v0.8.6) - 2023-08-30

### Added
1. Added optional --card-format command line argument to `put-acl` and `load-acl`
2. Implemented activate-keypads command to activate/deactivate reader access keypads.


## [0.8.5](https://github.com/uhppoted/uhppote-cli/releases/tag/v0.8.5) - 2023-06-13

### Added
1. `set-interlock` command.


## [0.8.4](https://github.com/uhppoted/uhppote-cli/releases/tag/v0.8.4) - 2023-03-17

### Added
1. `doc.go` package overview documentation.
2. Added (optional) PIN to put-card command (cf. https://github.com/uhppoted/uhppoted/issues/23)
3. Added --with-pin command line option to `get-acl`, `compare-acl` and `load-acl` (cf. https://github.com/uhppoted/uhppoted/issues/23)

### Updated
1. Included static-check in CI build.
2. Implemented custom formatting for `get-cards` to pretty print lists with too large card numbers.


## [0.8.3](https://github.com/uhppoted/uhppote-cli/releases/tag/v0.8.3) - 2022-12-16

### Added
1. Added ARM64 to release builds

### Changed
1. Removed _zip_ files from release artifacts (no longer necessary)


## [0.8.2](https://github.com/uhppoted/uhppote-cli/releases/tag/v0.8.2) - 2022-10-14

### Changed
1. Maintenance release for compatiblity with [uhppote-core](https://github.com/uhppoted/uhppote-core) v0.8.2
2. Bumped Go version to 1.19


## [0.8.1](https://github.com/uhppoted/uhppote-cli/releases/tag/v0.8.1) - 2022-08-01

### Changed
1. Maintenance release for compatiblity with [uhppote-core](https://github.com/uhppoted/uhppote-core) v0.8.1


## [0.8.0](https://github.com/uhppoted/uhppote-cli/releases/tag/v0.8.0) - 2022-07-01

### Changed
1. Maintenance release for compatiblity with [uhppote-core](https://github.com/uhppoted/uhppote-core) v0.8.0


### Changed
1. Maintenance release for compatibility with [uhppote-core](https://github.com/uhppoted/uhppote-core) 
   v0.8.0


## [0.7.3](https://github.com/uhppoted/uhppote-cli/releases/tag/v0.7.3) - 2022-06-01

### Changed
1. Maintenance release for compatibility with [uhppote-core](https://github.com/uhppoted/uhppote-core) 
   v0.7.3
2. Fixed errors in README (cf. https://github.com/uhppoted/uhppote-cli/issues/8)


## [0.7.2](https://github.com/uhppoted/uhppote-cli/releases/tag/v0.7.2) - 2022-01-26

### Changed
1. Updated command line arguments to support types.BindAddr, types.BroadcastAddr and types.ListenAddr
2. Replaced event rollover with handling for _nil_ and _overwritten_ events
3. Reworked `get-events` to also retrieve the current event index
4. Reworked `get-event` to support retrieving multiple events with _next:N_


## [0.7.1](https://github.com/uhppoted/uhppote-cli/releases/tag/v0.7.1) - 2021-07-01

### Added
1. Task lists:
   -  `clear-task-list`
   -  `add-task`
   -  `refresh-task-list`
   -  `set-task-list`

2. Documented commands missing from the README.


## [0.7.0](https://github.com/uhppoted/uhppote-cli/releases/tag/v0.7.0) - 2021-06-21

### Added
1. Support for using controller names (from uhppoted.conf) in device commands. e.g. 
   ```
   uhppote-cli get-device alpha
   ```

2. Time profiles:
   -  `get-time-profile`
   -  `set-time-profile`
   -  `clear-time-profiles`
   -  `get-time-profiles`
   -  `set-time-profiles`

3. Documented commands missing from the README.


## [0.6.12](https://github.com/uhppoted/uhppote-cli/releases/tag/v0.6.12) - 2021-04-28

### Changed
1. Added validation for `bind`, `broadcast` and `listen` ports to mitigate common misconfigurations.
2. Corrected typo in _README_

## Older

| *Version* | *Description*                                                                                      |
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
