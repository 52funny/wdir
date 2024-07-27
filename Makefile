.PHONY: build

# Development
build_rev := "main"
ifneq ($(wildcard .git),)
	build_rev := $(shell git rev-parse --short HEAD)
endif
build_date := $(shell date -u '+%Y-%m-%dT%H:%M:%S')

run: build
	@./wdir

build:
	@go build -ldflags "-s -w -X main.commit=$(build_rev) -X main.date=$(build_date)" -o wdir

install:
	@go install -ldflags "-s -w -X main.commit=$(build_rev) -X main.date=$(build_date)" .
