# TODO

- [x] [`restore-default-parameters`](https://github.com/uhppoted/uhppoted/issues/48)
      - [x] command
      - [x] CHANGELOG
      - [x] README

- [ ] CLI is waiting for CR on error

## TODO

- [ ] Rework command line parsing with tree-sitter
- [ ] Glamour
      - https://github.com/charmbracelet/glamour
- [ ] HOWTO: ACL with Google Sheets
      - `curl -Lo ACL.tsv "https://docs.google.com/spreadsheets/d/1_erZMyFmO6PM0PrAfEqdsiH9haiw-2UqY0kLwo_WTO8/export?gid=640947601&format=tsv"`
      - https://stackoverflow.com/questions/24255472/download-export-public-google-spreadsheet-as-tsv-from-command-line

- [ ] Windmill a la gcloud ...⠏⠹ (etc) 
- [ ] Unit/integration test for door control
- [ ] Restructure main()
      - https://pace.dev/blog/2020/02/12/why-you-shouldnt-use-func-main-in-golang-by-mat-ryer.html
- [ ] --changelog
      - https://bhupesh-v.github.io/why-how-add-changelog-in-your-next-cli/
- [ ] https://capiche.com/e/consumer-dev-tools-command-palette
- [ ] Check card number field for get-event
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
