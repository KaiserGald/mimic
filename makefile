# makefile for mimic
# 19 January 2018
# Code is licensed under the MIT License
# © 2018 Scott Isenberg

BINARY_NAME=mimic
BIN=bin/$(BINARY_NAME)
PSID:=$(shell pgrep $(BINARY_NAME))
OLD_INSTALL=$(GOBIN)/$(BINARY_NAME)
DONE=@echo -e $(GREEN)Done.$(NC)
RED='\033[0;31m'
GREEN='\033[0;32m'
WHITE='\033[1;37m'
PURPLE='\e[95m'
CYAN='\e[36m'
YELLOW='\033[1;33m'
ORANGE='\033[38;5;208m'
NC='\033[0m'
SED_COLORED=sed ''/'\(--- PASS\)'/s//$$(printf $(GREEN)---\\x20PASS)/'' | sed ''/PASS/s//$$(printf $(GREEN)PASS)/'' | sed  ''/'\(=== RUN\)'/s//$$(printf $(YELLOW)===\\x20RUN)/'' | sed ''/ok/s//$$(printf $(GREEN)ok)/'' | sed  ''/'\(--- FAIL\)'/s//$$(printf $(RED)---\\x20FAIL)/'' | sed  ''/FAIL/s//$$(printf $(RED)FAIL)/'' | sed ''/RUN/s//$$(printf $(YELLOW)RUN)/'' | sed ''/?/s//$$(printf $(ORANGE)?)/'' | sed ''/'\(^\)'/s//$$(printf $(NC))/''
ISSERVICERUNNING=$(shell pgrep $(BINARY_NAME))

all : stop deps test build install clean

build:
	@echo -e Building $(PURPLE)$(BINARY_NAME)$(NC)...
	@go build -o $(BIN) -v
	$(DONE)

clean:
	@echo -e Cleaning up...
	@go clean
	$(DONE)

deps:
	@echo -e Grabbing dependencies...
	@go get github.com/radovskyb/watcher
	$(DONE)

install:
	@echo -e Installing $(PURPLE)$(BINARY_NAME)$(NC) into $(CYAN)$(GOBIN)$(NC)...
	@echo Removing old install...
	@rm -f $(OLD_INSTALL)
	@echo Copying files... $(GOBING)
	@cp -u $(BIN) $(GOBIN)
	@cp -u $(BIN) /usr/local/bin
	$(DONE)

test:
	@echo -e Running Tests...
	@go test -args -w "testsrc:testdes" | ${SED_COLORED}
	@go test ./filewatcher/ | ${SED_COLORED}
	@go test ./filehandler/ | ${SED_COLORED}
	$(DONE)

run: all
	$(BINARY_NAME)

stop:
	@echo -e Checking if $(PURPLE)$(BINARY_NAME)$(NC) is running...
ifneq (${ISSERVICERUNNING},)
	@echo -e $(PURPLE)$(BINARY_NAME)$(NC) is running. Stopping it now.
	@kill $(PSID)
	$(DONE)
else
	@echo -e $(PURPLE)$(BINARY_NAME)$(NC) isn\'t currently running.
endif


.PHONY: all clean install
