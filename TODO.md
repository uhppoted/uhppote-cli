## v0.6x

**IN PROGRESS**

- [x] Rework load-acl to use `uhppoted-api`
- [x] get-acl
- [x] Rework `uhppoted-api::MakeTSV` to use `MakeTable`
- [x] compare-acl
- [x] Default to commmon config file
- [ ] Update documentation for ACL
- [ ] Rework grant/revoke for labelled doors
- [ ] Rework to use IDevice
- [ ] Install documentation
- [ ] Release v0.6.1

## TODO

### CLI
- [ ] Human readable output for e.g. get-status
- [ ] JSON formatted output for e.g. get-status
- [ ] Interactive shell (https://drewdevault.com/2019/09/02/Interactive-SSH-programs.html)
- [ ] use flag.FlagSet for commands
- [ ] Use (loadable) text/template for output formats
- [ ] Rework GetDevices to also find 'known' devices
- [ ] events: retrieve and show actual events
- [ ] Generate OTP secret + QR code

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
