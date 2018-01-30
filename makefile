# makefile for mimic
# 19 January 2018
# Code is licensed under the MIT License
# © 2018 Scott Isenberg

BINARY_NAME=mimic
BIN=bin/$(BINARY_NAME)
OLD_INSTALL=$(GOBIN)/$(BINARY_NAME)
DONE=@echo Done.

all : deps test build install clean

build:
	@echo Building $(BINARY_NAME)...
	@go build -o $(BIN) -v
	$(DONE)

clean:
	@echo Cleaning up...
	@go clean
	$(DONE)

deps:
	@echo Grabbing dependencies...
	@go get github.com/radovskyb/watcher
	$(DONE)

install:
	@echo Installing $(BINARY_NAME) into $(GOBIN)...
	@echo Removing old install...
	@rm -f $(OLD_INSTALL)
	@echo Copying files...
	@cp -u $(BIN) $(GOBIN)
	@sudo cp -u $(BIN) /usr/local/bin
	$(DONE)

test:
	@echo Running Tests...
	@go test -cover ./...
	$(DONE)

run: all
	$(BINARY_NAME)

.PHONY: all clean install
