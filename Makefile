VERSION = v0.7.x
LDFLAGS = -ldflags "-X uhppote.VERSION=$(VERSION)" 
DIST   ?= development
CLI     = ./bin/uhppote-cli

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
	go get -u github.com/uhppoted/uhppoted-api

debug: build
#	$(CLI) set-time-profile 405419896 2  2021-01-01:2021-12-31 Mon,Thurs,Sat 09:30-12:30,13:45-16:00,19:30-20:30 
#	$(CLI) set-time-profile 405419896 29 2021-04-01:2021-12-31 Mon,Wed,Fri 08:30-11:30,,13:45-17:00 2
#	$(CLI) set-time-profile 405419896 55 2021-04-01:2021-10-31 Sat,Sun 10:30-11:30 
#	$(CLI) get-time-profiles $(SERIALNO) 
	$(CLI) get-time-profiles $(SERIALNO) ../runtime/profiles.tsv

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

set-address: build
	$(CLI) $(DEBUG) set-address $(SERIALNO) $(DEVICEIP) '255.255.255.0' '0.0.0.0'

get-status: build
	$(CLI) $(DEBUG) get-status $(SERIALNO)

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

get-listener: build
	$(CLI) $(DEBUG) get-listener $(SERIALNO)

set-listener: build
	$(CLI) $(DEBUG) set-listener $(SERIALNO) $(LISTEN)

get-cards: build
	$(CLI) $(DEBUG) get-cards $(SERIALNO)

get-card: build
	$(CLI) $(DEBUG) get-card $(SERIALNO) $(CARD)

put-card: build
	$(CLI) $(DEBUG) put-card $(SERIALNO) $(CARD) 2021-01-01 2021-12-31 1,3,4:29

delete-card: build
	$(CLI) $(DEBUG) delete-card $(SERIALNO) $(CARD)

delete-all: build
#	$(CLI) $(DEBUG) delete-all $(SERIALNO)
	$(CLI) $(DEBUG) delete-all 405419896
	$(CLI) $(DEBUG) delete-all 303986753

get-time-profile: build
	$(CLI) --debug get-time-profile $(SERIALNO)   29
	# $(CLI) --debug get-time-profile $(CONTROLLER) 2
	# $(CLI) --debug get-time-profile 423187757   29

get-time-profiles: build
	$(CLI) get-time-profiles $(SERIALNO) 
	$(CLI) get-time-profiles $(SERIALNO) ../runtime/profiles.tsv

set-time-profile: build
	$(CLI) --debug set-time-profile 303986753 29 2021-04-01:2021-12-31 Mon,Wed,Fri 08:30-11:30,,13:45-17:00 
	$(CLI) --debug set-time-profile 405419896 29 2021-04-01:2021-12-31 Mon,Wed,Fri 08:30-11:30,,13:45-17:00 
	# $(CLI) --debug set-time-profile $(CONTROLLER) 2 2021-04-01:2021-12-31 Mon,Wed,Fri 08:30-11:30,,13:45-17:00 3
	# $(CLI) --debug set-time-profile 423187757 29 2021-04-01:2021-12-31 Mon,Wed,Fri 08:30-11:30,,13:45-17:00 

clear-time-profiles: build
	$(CLI) --debug clear-time-profiles $(SERIALNO)
	# $(CLI) --debug clear-time-profiles $(CONTROLLER)
	# $(CLI) --debug clear-time-profiles 423187757

get-events: build
	$(CLI) $(DEBUG) get-events $(SERIALNO)

get-event: build
#	$(CLI) get-event $(SERIALNO) 17
#	$(CLI) get-event $(SERIALNO) 17263
#	$(CLI) get-event $(SERIALNO) first
#	$(CLI) get-event $(SERIALNO) last
	$(CLI) get-event $(SERIALNO)

get-event-index: build
	$(CLI) $(DEBUG) get-event-index $(SERIALNO)

set-event-index: build
	$(CLI) $(DEBUG) set-event-index $(SERIALNO) 23

record-special-events: build
	$(CLI) $(DEBUG) record-special-events $(SERIALNO) true

open: build
	$(CLI) $(DEBUG) open $(SERIALNO) 1

listen: build
	$(CLI) --listen $(LISTEN) $(DEBUG) listen 

# ACL COMMANDS

show: build
	$(CLI) $(DEBUG) show $(CARD)

grant: build
	$(CLI) $(DEBUG) grant $(CARD) 2020-01-01 2020-12-31 "Lady's Chamber, Workshop"

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

