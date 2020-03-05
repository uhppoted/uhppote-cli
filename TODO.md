## v0.6.0x

**IN PROGRESS**

## TODO

### CLI
- [ ] Rework grant/revoke for individual doors (labelled)
- [ ] get-acl
- [ ] Human readable output for e.g. get-status
- [ ] JSON formatted output for e.g. get-status
- [ ] Interactive shell (https://drewdevault.com/2019/09/02/Interactive-SSH-programs.html)
- [ ] use flag.FlagSet for commands
- [ ] Default to commmon config file
- [ ] Use (loadable) text/template for output formats
- [ ] Rework GetDevices to also find 'known' devices
- [ ] events: retrieve and show actual events
- [ ] Generate OTP secret + QR code

### Documentation

- [ ] godoc
- [ ] build documentation
- [ ] install documentation
- [ ] user manuals
- [ ] man/info page

### Other

1.  Update to use modules
2.  Consistently include device serial number in output e.g. of get-time
3.  github project page
4.  Integration tests
5.  Verify fields in listen events/status replies against SDK:
    - battery status can be (at least) 0x00, 0x01 and 0x04
