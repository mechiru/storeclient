SHELL := /bin/bash

CURR_DIR := $(shell pwd)
MODULES := $(shell go work edit -json .go.work | jq ".Use[].DiskPath")
GO_VERSION := 1.17

.PHONY: tidy
tidy:
	@for dir in $(MODULES); do \
		cd "$(CURR_DIR)/$$dir" \
		&& mod=$$(go mod edit -json | jq -r .Module.Path) \
		&& rm -f go.{mod,sum} \
		&& go mod init $$mod \
		&& go mod edit -go $(GO_VERSION) \
		&& go mod tidy; \
	done

.PHONY: build
build:
	@for dir in $(MODULES); do \
		if [ $$dir = "." ]; then continue; fi; \
		cd "$(CURR_DIR)/$$dir" && echo "build $$(pwd)..." && go build -v ./...; \
	done

.PHONY: test
test:
	@for dir in $(MODULES); do \
		if [ $$dir = "." ]; then continue; fi; \
		cd "$(CURR_DIR)/$$dir" && echo "test $$(pwd)..." && go test -v ./...; \
	done
