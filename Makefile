ALL_GO_MOD_DIRS := $(shell find . -mindepth 2 -type f -name 'go.mod' -exec dirname {} \; | sort)

.PHONY: build
build:
	@set -e; \
	for dir in $(ALL_GO_MOD_DIRS); do \
		echo "go build in $${dir}"; \
		(cd "$${dir}" && go build -v); \
	done;

.PHONY: test
test:
	@set -e; \
	for dir in $(ALL_GO_MOD_DIRS); do \
		echo "go test ./... in $${dir}"; \
		(cd "$${dir}" && go test -v ./...); \
	done;
