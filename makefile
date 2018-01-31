# makefile for mimic
# 19 January 2018
# Code is licensed under the MIT License
# © 2018 Scott Isenberg

BINARY_NAME=mimic
BIN=bin/$(BINARY_NAME)
OLD_INSTALL=$(GOBIN)/$(BINARY_NAME)
DONE=@echo -e $(GREEN)Done.$(NC)
RED='\033[0;31m'
GREEN='\033[0;32m'
WHITE='\033[1;37m'
PURPLE='\e[95m'
CYAN='\e[36m'
YELLOW='\033[1;33m'
NC='\033[0m'

all : deps test build install clean

build:
	@echo -e $(WHITE)Building $(PURPLE)$(BINARY_NAME)$(WHITE)...
	@go build -o $(BIN) -v
	$(DONE)

clean:
	@echo -e $(WHITE)Cleaning up...
	@go clean
	$(DONE)

deps:
	@echo -e $(WHITE)Grabbing dependencies...
	@go get github.com/radovskyb/watcher
	$(DONE)

install:
	@echo -e $(WHITE)Installing $(PURPLE)$(BINARY_NAME)$(WHITE) into $(CYAN)$(GOBIN)$(WHITE)...
	@echo Removing old install...
	@rm -f $(OLD_INSTALL)
	@echo Copying files...
	@cp -u $(BIN) $(GOBIN)
	@sudo cp -u $(BIN) /usr/local/bin
	$(DONE)

test:
	@echo -e $(WHITE)Running Tests...
	@go test ./... | sed ''/'\(--- PASS\)'/s//$$(printf $(GREEN)---\\x20PASS$(WHITE))/'' | sed ''/PASS/s//$$(printf $(GREEN)PASS$(WHITE))/'' | sed  ''/'\(=== RUN\)'/s//$$(printf $(YELLOW)===\\x20RUN$(WHITE))/'' | sed ''/ok/s//$$(printf $(GREEN)ok$(WHITE))/'' | sed  ''/'\(--- FAIL\)'/s//$$(printf $(RED)---\\x20FAIL$(WHITE))/'' | sed  ''/FAIL/s//$$(printf $(RED)FAIL$(WHITE))/'' | sed ''/RUN/s//$$(printf $(YELLOW)RUN$(WHITE))/''
	$(DONE)

run: all
	$(BINARY_NAME)

.PHONY: all clean install
