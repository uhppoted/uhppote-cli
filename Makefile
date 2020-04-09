VERSION = v0.6.0x
LDFLAGS = -ldflags "-X uhppote.VERSION=$(VERSION)" 
DIST   ?= development
CLI     = ./bin/uhppote-cli

SERIALNO ?= 405419896
CARD     ?= 65538
DOOR     ?= 3
DEVICEIP ?= 192.168.1.125
DATETIME  = $(shell date "+%Y-%m-%d %H:%M:%S")
LISTEN   ?= 192.168.1.100:60001
DEBUG    ?= --debug

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

release: test vet
	mkdir -p dist/$(DIST)/windows
	mkdir -p dist/$(DIST)/darwin
	mkdir -p dist/$(DIST)/linux
	mkdir -p dist/$(DIST)/arm7
	env GOOS=linux   GOARCH=amd64         go build -o dist/$(DIST)/linux   ./...
	env GOOS=linux   GOARCH=arm   GOARM=7 go build -o dist/$(DIST)/arm7    ./...
	env GOOS=darwin  GOARCH=amd64         go build -o dist/$(DIST)/darwin  ./...
	env GOOS=windows GOARCH=amd64         go build -o dist/$(DIST)/windows ./...

release-tar: release
	find . -name ".DS_Store" -delete
	tar --directory=dist --exclude=".DS_Store" -cvzf dist/$(DIST).tar.gz $(DIST)
	cd dist; zip --recurse-paths $(DIST).zip $(DIST)

debug: build
	$(CLI) revoke $(SERIALNO) $(CARD)
	$(CLI) get-card $(SERIALNO) $(CARD)
	$(CLI) put-card $(SERIALNO) $(CARD) 2020-01-01 2020-12-31 1,4
# #	$(CLI) grant $(SERIALNO) $(CARD) 2020-01-01 2020-12-31 1,4
# 	$(CLI) grant $(CARD) 2020-01-01 2020-12-31 Workshop
	$(CLI) get-card $(SERIALNO) $(CARD)
	$(CLI) put-card $(SERIALNO) $(CARD) 2020-01-01 2020-12-31 2,3
	$(CLI) get-card $(SERIALNO) $(CARD)

usage: build
	$(CLI)

help: build
	$(CLI) help
	$(CLI) help get-devices

version: build
	$(CLI) version

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
	$(CLI) $(DEBUG) put-card $(SERIALNO) $(CARD) 2020-01-01 2020-12-31 1,2,3,4

revoke: build
	$(CLI) $(DEBUG) revoke $(SERIALNO) $(CARD)

revoke-all: build
	$(CLI) $(DEBUG) revoke-all $(SERIALNO)

load-acl: build
	$(CLI) $(DEBUG) --config ../runtime/simulation/$(SERIALNO).conf load-acl ../runtime/simulation/$(SERIALNO).acl

get-acl: build
	$(CLI) get-acl
	$(CLI) $(DEBUG) --config ../runtime/simulation/$(SERIALNO).conf get-acl ../runtime/simulation/uhppote-cli.acl

compare-acl: build
	$(CLI) $(DEBUG) compare-acl ../runtime/simulation/simulation.acl
	$(CLI) $(DEBUG) --config ../runtime/simulation/$(SERIALNO).conf compare-acl ../runtime/simulation/simulation.acl ../runtime/simulation/$(SERIALNO).rpt

get-events: build
	$(CLI) $(DEBUG) get-events $(SERIALNO)

get-event-index: build
	$(CLI) $(DEBUG) get-event-index $(SERIALNO)

set-event-index: build
	$(CLI) $(DEBUG) set-event-index $(SERIALNO) 23

open: build
	$(CLI) $(DEBUG) open $(SERIALNO) 1

listen: build
	$(CLI) --listen $(LISTEN) $(DEBUG) listen 

