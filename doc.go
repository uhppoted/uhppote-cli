// Copyright 2023 uhppoted@twyst.co.za. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

/*
Package uhppote-cli implements a command line interface to the UHPPOTE TCP/IP Wiegand-26 access controllers.

The CLI provides a basic set of commands to:

  - query and configure the access controller
  - manage access control lists
  - receive and fetch events

Supported commands:
  - get-devices
  - get-device
  - set-address
  - get-time
  - set-time
  - get-door-delay
  - set-door-delay
  - get-door-control
  - set-door-control
  - get-listener
  - set-listener
  - record-special-events
  - get-status
  - get-cards
  - get-card
  - put-card
  - delete-card
  - delete-all
  - get-time-profile
  - set-time-profile
  - get-time-profiles
  - set-time-profiles
  - clear-time-profiles
  - clear-task-list
  - add-task
  - refresh-task-list
  - set-task-list
  - get-events
  - get-event
  - get-event-index
  - set-event-index
  - open
  - listen
  - grant
  - revoke
  - load-acl
  - get-acl
  - compare-acl
*/
package cli
