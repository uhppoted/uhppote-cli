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

.PHONY: bump

all: test      \
	 benchmark \
     coverage

clean:
	go clean
	rm -rf bin

format: 
	go fmt ./...

build: format
	mkdir -p bin
	go build -o bin ./...

test: build
	go test ./...

vet: build
	go vet ./...

lint: build
	golint ./...

benchmark: build
	go test -bench ./...

coverage: build
	go test -cover ./...

build-all: test vet
	mkdir -p dist/$(DIST)/windows
	mkdir -p dist/$(DIST)/darwin
	mkdir -p dist/$(DIST)/linux
	mkdir -p dist/$(DIST)/arm7
	env GOOS=linux   GOARCH=amd64         go build -o dist/$(DIST)/linux   ./...
	env GOOS=linux   GOARCH=arm   GOARM=7 go build -o dist/$(DIST)/arm7    ./...
	env GOOS=darwin  GOARCH=amd64         go build -o dist/$(DIST)/darwin  ./...
	env GOOS=windows GOARCH=amd64         go build -o dist/$(DIST)/windows ./...

release: build-all
	find . -name ".DS_Store" -delete
	tar --directory=dist --exclude=".DS_Store" -cvzf dist/$(DIST).tar.gz $(DIST)
	cd dist; zip --recurse-paths $(DIST).zip $(DIST)

bump:
	go get -u github.com/uhppoted/uhppote-core
	go get -u github.com/uhppoted/uhppoted-lib

debug: build
	$(CLI) version
	$(CLI) --config '../runtime/CLI/uhppoted.conf' version
	$(CLI) --bind '192.168.1.100:54321' --broadcast '192.168.1.100:43210' --listen '192.168.1.100:32109' version
	$(CLI) --config '../runtime/CLI/uhppoted.conf' --bind '192.168.1.100:54321' --broadcast '192.168.1.100:43210' --listen '192.168.1.100:32109' version

irl: build
	$(CLI) set-time            423187757
	$(CLI) clear-time-profiles 423187757
	$(CLI) set-time-profile    423187757 29 2021-01-01:2021-12-31 Mon 08:30-17:00
	$(CLI) get-time-profiles   423187757
	$(CLI) put-card            423187757 6154410 2021-01-01 2021-12-31 3:29
	$(CLI) clear-task-list     423187757
	$(CLI) add-task            423187757 'disable time profile' 3 2021-01-01:2021-12-31 Mon 08:30
	$(CLI) add-task            423187757 'enable time profile'  3 2021-01-01:2021-12-31 Mon 11:30
	$(CLI) add-task            423187757 'lock door'            3 2021-01-01:2021-12-31 Mon 11:45
	$(CLI) add-task            423187757 'unlock door'          3 2021-01-01:2021-12-31 Mon 12:00
	$(CLI) add-task            423187757 'control door'         3 2021-01-01:2021-12-31 Mon 12:05
	$(CLI) add-task            423187757 'disable pushbutton'   3 2021-01-01:2021-12-31 Mon 12:10
	$(CLI) add-task            423187757 'enable pushbutton'    3 2021-01-01:2021-12-31 Mon 12:15
	$(CLI) add-task            423187757 'trigger once'         3 2021-01-01:2021-12-31 Mon 12:05
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
	$(CLI) $(DEBUG) put-card $(SERIALNO) $(CARD) 2021-01-01 2021-12-31 1,3,4:29

delete-card: build
	$(CLI) $(DEBUG) delete-card $(SERIALNO) $(CARD)

delete-all: build
	$(CLI) delete-all $(SERIALNO)

get-time-profile: build
	$(CLI) $(DEBUG) get-time-profile $(SERIALNO) 29

set-time-profile: build
	$(CLI) $(DEBUG) set-time-profile 303986753 29 2021-04-01:2021-12-31 Mon,Wed,Fri 08:30-11:30,,13:45-17:00 
	$(CLI) $(DEBUG) set-time-profile 405419896 29 2021-04-01:2021-12-31 Mon,Wed,Fri 08:30-11:30,,13:45-17:00 3

clear-time-profiles: build
	$(CLI) $(DEBUG) clear-time-profiles $(SERIALNO)

get-time-profiles: build
	$(CLI) get-time-profiles $(SERIALNO) 
	$(CLI) get-time-profiles $(CONTROLLER) ../runtime/$(CONTROLLER).tsv

set-time-profiles: build
	$(CLI) clear-time-profiles $(SERIALNO) 
	$(CLI) set-time-profile $(SERIALNO) 75  2021-04-01:2021-12-31 Mon,Wed,Fri 08:30-11:30,,13:45-17:00 
	$(CLI) set-time-profile $(SERIALNO) 100 2021-04-01:2021-12-31 Mon,Wed,Fri 08:30-11:30,,13:45-17:00 
	$(CLI) set-time-profile 303986753   101 2021-01-01:2021-12-31 Sat,Sun     10:30-16:30
	$(CLI) set-time-profiles $(SERIALNO) ../runtime/set-time-profiles.tsv
	$(CLI) get-time-profiles $(SERIALNO) 

clear-task-list: build
	$(CLI) --debug clear-task-list $(SERIALNO)

add-task: build
	$(CLI) --debug add-task $(SERIALNO) 3 4 2021-01-01:2021-12-31 Mon,Fri 08:30 33
	$(CLI) --debug add-task $(SERIALNO) 'enable more cards' 4 2021-01-01:2021-12-31 Mon,Fri 08:30 29

refresh-task-list: build
	$(CLI) --debug refresh-task-list $(SERIALNO)

set-task-list: build
	$(CLI) set-task-list $(SERIALNO) ../runtime/set-tasks.tsv

get-events: build
	$(CLI) $(DEBUG) get-events $(SERIALNO)

get-event: build
	$(CLI) get-event $(SERIALNO) 17
	$(CLI) get-event $(SERIALNO) 17263
	$(CLI) get-event $(SERIALNO) first
	$(CLI) get-event $(SERIALNO) last
	$(CLI) get-event $(SERIALNO)

get-event-index: build
	$(CLI) $(DEBUG) get-event-index $(SERIALNO)

set-event-index: build
	$(CLI) $(DEBUG) set-event-index $(SERIALNO) 23

open: build
	$(CLI) $(DEBUG) open $(SERIALNO) 1

listen: build
	$(CLI) --listen $(LISTEN) $(DEBUG) listen 

# ACL COMMANDS

show: build
	$(CLI) show $(CARD)

grant: build
	$(CLI) grant $(CARD) 2021-01-01 2021-12-31 "Gryffindor, Slytherin"
	$(CLI) grant $(CARD) 2021-01-01 2021-12-31 29 "Dungeon"

grant-all: build
	$(CLI) $(DEBUG) grant $(CARD) 2020-01-01 2020-12-31 ALL

revoke: build
	$(CLI) $(DEBUG) revoke $(CARD) "Lady's Chamber, D2"

revoke-all: build
	$(CLI) $(DEBUG) revoke $(CARD) ALL
	
load-acl: build
	$(CLI) --config ../runtime/simulation/$(SERIALNO).conf load-acl ../runtime/simulation/$(SERIALNO).acl

get-acl: build
	$(CLI) get-acl
#	$(CLI) $(DEBUG) --config ../runtime/simulation/$(SERIALNO).conf get-acl ../runtime/simulation/uhppote-cli.acl

compare-acl: build
	$(CLI) $(DEBUG) compare-acl ../runtime/simulation/simulation.acl
	$(CLI) $(DEBUG) --config ../runtime/simulation/$(SERIALNO).conf compare-acl ../runtime/simulation/simulation.acl ../runtime/simulation/$(SERIALNO).rpt

