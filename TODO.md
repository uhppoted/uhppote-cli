## v0.7.x

### IN PROGRESS

- [ ] Tabulate output for get-devices
- [ ] Add missing commands to README
      - get-events
      - get-event-index
      - set-event-index
      - open
      - listen
- [ ] `set-schedule`
- [ ] `get-schedule`

- [x] Print device name for get-device(s)
- [x] Check for duplicate time profiles
- [x] Replace all (deprecated) ioutil.WriteFile with os.WriteFile
- [x] Update README with time profile commands
- [x] Verify 'to' is not before 'from'
- [x] Verify segment 'end' is not before segment 'start'
- [x] `set-time-profiles`
- [x] `get-time-profiles`
- [x] `compare-acl` with time profiles
- [x] `load-acl` with time profiles
- [x] `get-acl` with time profiles
- [x] Add time profiles to put-card
- [x] Check for linked circular references in set-time-profile
- [x] `clear-time-profiles`
- [x] `get-time-profile`
- [x] `set-time-profile`
- [x] Use device name from conf file

## TODO

- [ ] Route debugging to stderr
- [ ] get-events --fetch
- [ ] listener: retrieve and show actual events

- [ ] Progress messages for acl-load
- [ ] Nicer formatting for acl-xxx
- [ ] Human readable output for e.g. get-status
- [ ] JSON formatted output for e.g. get-status
- [ ] Interactive shell (https://drewdevault.com/2019/09/02/Interactive-SSH-programs.html)
- [ ] use flag.FlagSet for commands
- [ ] Use (loadable) text/template for output formats
- [ ] Generate OTP secret + QR code
- [ ] --no-log option to suppress progress messages

### Documentation

- [ ] godoc
- [ ] build documentation
- [ ] user manuals
- [ ] man/info page

### Other

1.  Consistently include device serial number in output e.g. of get-time
2.  Integration tests
3.  Verify fields in listen events/status replies against SDK:
    - battery status can be (at least) 0x00, 0x01 and 0x04

### Miscellaneous

1. [syncthing](https://tonsky.me/blog/syncthing/?utm_source=hackerbits.com&utm_medium=email&utm_campaign=issue54)
2. bash scripts to retrieve all events:
   ```
   -- get-event
   #/bin/bash
   ./bin/uhppote-cli get-event $1

   -- get-events
   #/bin/bash
   N=1
   while [ $n -le 5 ]
   do
      ./get-event 405419896
      N=$(( N+1 ))
   done

   ./get-events 1> >(tee -a x.log y.log 1> /dev/null) 2>> errors.log
   ```
