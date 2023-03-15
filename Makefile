DIST ?= development
CLI   = ./bin/uhppote-cli

CONTROLLER ?= Alpha
SERIALNO ?= 405419896
CARD     ?= 8165538
DOOR     ?= 3
DEVICEIP ?= 192.168.1.125
DATETIME  = $(shell date "+%Y-%m-%d %H:%M:%S")
LISTEN   ?= 192.168.1.100:60001
DEBUG    ?= --debug

.PHONY: clean
.PHONY: update
.PHONY: update-release

all: test      \
	 benchmark \
     coverage

clean:
	go clean
	rm -rf bin

update:
	go get -u github.com/uhppoted/uhppote-core@master
	go get -u github.com/uhppoted/uhppoted-lib@master
	go mod tidy

update-release:
	go get -u github.com/uhppoted/uhppote-core
	go get -u github.com/uhppoted/uhppoted-lib
	go mod tidy

format: 
	go fmt ./...

build: format
	mkdir -p bin
	go build -trimpath -o bin ./...

test: build
	go test ./...

vet: 
	go vet ./...

lint: 
	env GOOS=darwin  GOARCH=amd64 staticcheck ./...
	env GOOS=linux   GOARCH=amd64 staticcheck ./...
	env GOOS=windows GOARCH=amd64 staticcheck ./...

benchmark: build
	go test -bench ./...

coverage: build
	go test -cover ./...

build-all: test vet lint
	mkdir -p dist/$(DIST)/windows
	mkdir -p dist/$(DIST)/darwin
	mkdir -p dist/$(DIST)/linux
	mkdir -p dist/$(DIST)/arm
	mkdir -p dist/$(DIST)/arm7
	env GOOS=linux   GOARCH=amd64         GOWORK=off go build -trimpath -o dist/$(DIST)/linux   ./...
	env GOOS=linux   GOARCH=arm64         GOWORK=off go build -trimpath -o dist/$(DIST)/arm     ./...
	env GOOS=linux   GOARCH=arm   GOARM=7 GOWORK=off go build -trimpath -o dist/$(DIST)/arm7    ./...
	env GOOS=darwin  GOARCH=amd64         GOWORK=off go build -trimpath -o dist/$(DIST)/darwin  ./...
	env GOOS=windows GOARCH=amd64         GOWORK=off go build -trimpath -o dist/$(DIST)/windows ./...

release: update-release build-all
	find . -name ".DS_Store" -delete
	tar --directory=dist --exclude=".DS_Store" -cvzf dist/$(DIST).tar.gz $(DIST)
	cd dist; zip --recurse-paths $(DIST).zip $(DIST)

publish: release
	echo "Releasing version $(VERSION)"
	gh release create "$(VERSION)" "./dist/uhppote-cli_$(VERSION).tar.gz" "./dist/uhppote-cli_$(VERSION)*.zip" --draft --prerelease --title "$(VERSION)-beta" --notes-file release-notes.md

debug: build
	$(CLI) $(DEBUG) set-time-profile 405419896 3  2023-01-01:2023-12-31 Sat,Sun     09:30-16:30,, 
	$(CLI) $(DEBUG) set-time-profile 405419896 29 2023-04-01:2023-12-31 Mon,Wed,Fri 08:30-11:30,,13:45-17:00 3
	$(CLI) $(DEBUG) set-time-profile 303986753 3  2023-01-01:2023-12-31 Sat,Sun     09:30-16:30,, 
	$(CLI) $(DEBUG) set-time-profile 303986753 29 2023-04-01:2023-12-31 Mon,Wed,Fri 08:30-11:30,,13:45-17:00 3

irl: build
	$(CLI) set-time            423187757
	$(CLI) clear-time-profiles 423187757
	$(CLI) set-time-profile    423187757 29 2023-01-01:2023-12-31 Mon 08:30-17:00
	$(CLI) get-time-profiles   423187757
	$(CLI) put-card            423187757 6154410 2023-01-01 2023-12-31 3:29
	$(CLI) clear-task-list     423187757
	$(CLI) add-task            423187757 'disable time profile' 3 2023-01-01:2023-12-31 Mon 08:30
	$(CLI) add-task            423187757 'enable time profile'  3 2023-01-01:2023-12-31 Mon 11:30
	$(CLI) add-task            423187757 'lock door'            3 2023-01-01:2023-12-31 Mon 11:45
	$(CLI) add-task            423187757 'unlock door'          3 2023-01-01:2023-12-31 Mon 12:00
	$(CLI) add-task            423187757 'control door'         3 2023-01-01:2023-12-31 Mon 12:05
	$(CLI) add-task            423187757 'disable pushbutton'   3 2023-01-01:2023-12-31 Mon 12:10
	$(CLI) add-task            423187757 'enable pushbutton'    3 2023-01-01:2023-12-31 Mon 12:15
	$(CLI) add-task            423187757 'trigger once'         3 2023-01-01:2023-12-31 Mon 12:05
	$(CLI) refresh-task-list   423187757

godoc:
	godoc -http=:80	-index_interval=60s

usage: build
	$(CLI)

help: build
	$(CLI) help
	$(CLI) help get-devices
	$(CLI) help record-special-events

version: build
	$(CLI) version

# DEVICE COMMANDS

get-devices: build
	$(CLI) $(DEBUG) get-devices

get-device: build
	$(CLI) $(DEBUG) get-device $(SERIALNO)
	$(CLI) get-device 303986753
	$(CLI) get-device 405419896

set-address: build
	$(CLI) $(DEBUG) set-address $(SERIALNO) $(DEVICEIP) '255.255.255.0' '0.0.0.0'

get-listener: build
	$(CLI) $(DEBUG) get-listener $(SERIALNO)

set-listener: build
	$(CLI) $(DEBUG) set-listener $(SERIALNO) $(LISTEN)

get-time: build
	$(CLI) $(DEBUG) get-time $(SERIALNO)

set-time: build
	$(CLI) $(DEBUG) set-time $(SERIALNO)
	$(CLI) $(DEBUG) set-time $(SERIALNO) "$(DATETIME)"

get-door-delay: build
	$(CLI) $(DEBUG) get-door-delay $(SERIALNO) $(DOOR)

set-door-delay: build
	$(CLI) $(DEBUG) set-door-delay $(SERIALNO) $(DOOR) 5

get-door-control: build
	$(CLI) $(DEBUG) get-door-control $(SERIALNO) $(DOOR)

set-door-control: build
	$(CLI) $(DEBUG) set-door-control $(SERIALNO) $(DOOR) 'normally closed'

record-special-events: build
	$(CLI) $(DEBUG) record-special-events $(SERIALNO) true

get-status: build
	$(CLI) $(DEBUG) get-status $(SERIALNO)

get-cards: build
	$(CLI) $(DEBUG) get-cards $(SERIALNO)

get-card: build
	$(CLI) $(DEBUG) get-card $(SERIALNO) $(CARD)

put-card: build
	$(CLI) $(DEBUG) put-card $(SERIALNO) $(CARD) 2023-01-01 2023-12-31 1,3,4:29 7531

delete-card: build
	$(CLI) $(DEBUG) delete-card $(SERIALNO) $(CARD)

delete-cards: build
	$(CLI) delete-all $(SERIALNO)

get-time-profile: build
	$(CLI) $(DEBUG) get-time-profile $(SERIALNO) 29

set-time-profile: build
	$(CLI) $(DEBUG) set-time-profile 405419896 29 2023-04-01:2023-12-31 Mon,Wed,Fri 08:30-11:30,,13:45-17:00 
	$(CLI) $(DEBUG) set-time-profile 405419896 3  2023-01-01:2023-12-31 Sat,Sun     09:30-16:30,, 
	$(CLI) $(DEBUG) set-time-profile 405419896 29 2023-04-01:2023-12-31 Mon,Wed,Fri 08:30-11:30,,13:45-17:00 3

clear-time-profiles: build
	$(CLI) $(DEBUG) clear-time-profiles $(SERIALNO)

get-time-profiles: build
	$(CLI) get-time-profiles $(SERIALNO) 
	$(CLI) get-time-profiles $(CONTROLLER) ../runtime/$(CONTROLLER).tsv

set-time-profiles: build
	$(CLI) clear-time-profiles $(SERIALNO) 
	$(CLI) set-time-profile $(SERIALNO) 75  2023-04-01:2023-12-31 Mon,Wed,Fri 08:30-11:30,,13:45-17:00 
	$(CLI) set-time-profile $(SERIALNO) 100 2023-04-01:2023-12-31 Mon,Wed,Fri 08:30-11:30,,13:45-17:00 
	$(CLI) set-time-profile 303986753   101 2023-01-01:2023-12-31 Sat,Sun     10:30-16:30
	$(CLI) set-time-profiles $(SERIALNO) ../runtime/set-time-profiles.tsv
	$(CLI) get-time-profiles $(SERIALNO) 

clear-task-list: build
	$(CLI) --debug clear-task-list $(SERIALNO)

add-task: build
	$(CLI) --debug add-task $(SERIALNO) 3 4 2023-01-01:3-12-31 Mon,Fri 08:30 33
	$(CLI) --debug add-task $(SERIALNO) 'enable more cards' 4 2023-01-01:2023-12-31 Mon,Fri 08:30 29

refresh-task-list: build
	$(CLI) --debug refresh-task-list $(SERIALNO)

set-task-list: build
	$(CLI) set-task-list $(SERIALNO) ../runtime/set-tasks.tsv

get-events: build
	$(CLI) $(DEBUG) get-events $(SERIALNO)

get-event: build
	$(CLI) get-event $(SERIALNO) 17
	$(CLI) get-event $(SERIALNO) first
	$(CLI) get-event $(SERIALNO) last
	$(CLI) get-event $(SERIALNO) current
	$(CLI) get-event $(SERIALNO) next
	$(CLI) get-event $(SERIALNO) next:5
	$(CLI) get-event $(SERIALNO)
	$(CLI) get-event $(SERIALNO) 17263
	# $(CLI) get-event 201020304 100

get-event-index: build
	$(CLI) $(DEBUG) get-event-index $(SERIALNO)

set-event-index: build
	$(CLI) $(DEBUG) set-event-index $(SERIALNO) 23

open: build
	$(CLI) $(DEBUG) open $(SERIALNO) 1

set-pc-control: build
	$(CLI) $(DEBUG) set-pc-control $(SERIALNO) true
	# $(CLI) $(DEBUG) set-pc-control 423187757 true

listen: build
	$(CLI) --listen $(LISTEN) $(DEBUG) listen 

# ACL COMMANDS

show: build
	$(CLI) show $(CARD)

grant: build
	$(CLI) grant $(CARD) 2023-01-01 2023-12-31 "Gryffindor, Slytherin"
	$(CLI) grant $(CARD) 2023-01-01 2023-12-31 29 "Dungeon"

grant-all: build
	$(CLI) $(DEBUG) grant $(CARD) 2023-01-01 2023-12-31 ALL

revoke: build
	$(CLI) $(DEBUG) revoke $(CARD) "Lady's Chamber, D2"

revoke-all: build
	$(CLI) $(DEBUG) revoke $(CARD) ALL
	
get-acl: build
	$(CLI) get-acl
#	$(CLI) $(DEBUG) --config ../runtime/simulation/$(SERIALNO).conf get-acl ../runtime/simulation/uhppote-cli.acl

get-acl-with-pin: build
	$(CLI) get-acl
	$(CLI) get-acl --with-pin
	$(CLI) get-acl --with-pin ../runtime/uhppote-cli/acl-with-pin.tsv

compare-acl: build
	$(CLI) $(DEBUG) compare-acl ../runtime/simulation/simulation.acl
	$(CLI) $(DEBUG) --config ../runtime/simulation/$(SERIALNO).conf compare-acl ../runtime/simulation/simulation.acl ../runtime/simulation/$(SERIALNO).rpt

compare-acl-with-pin: build
	$(CLI) compare-acl ../runtime/simulation/simulation.acl
	$(CLI) compare-acl ../runtime/simulation/simulation.acl ../runtime/uhppote-cli/compare-acl.tsv
	$(CLI) compare-acl --with-pin ../runtime/simulation/simulation.acl
	$(CLI) compare-acl --with-pin ../runtime/simulation/simulation-with-pin.acl
	$(CLI) compare-acl --with-pin ../runtime/simulation/simulation.acl ../runtime/uhppote-cli/compare-acl-with-pin.tsv

load-acl: build
	$(CLI) --config ../runtime/simulation/$(SERIALNO).conf load-acl ../runtime/simulation/$(SERIALNO).acl

load-acl-with-pin: build
	$(CLI) load-acl ../runtime/simulation/simulation.acl
	$(CLI) load-acl --with-pin ../runtime/simulation/simulation-with-pin.acl

