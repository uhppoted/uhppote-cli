# CHANGELOG

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
