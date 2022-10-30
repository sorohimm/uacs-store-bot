# suppress output, run `make XXX V=` to be verbose
V := @

RELEASE=$(shell git describe --always --tags)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

OUT_DIR := ./build


default: all

.PHIONY: all
all:bot

.PHONY: linux
linux: export GOOS := linux
linux: export GOARCH := amd64
linux:all

#### api service build
.PHONY: bot
bot: BOT_OUT := $(OUT_DIR)/bot
bot: BOT_MAIN := ./cmd/bot
bot:
	@echo BUILDING $(BOT_OUT)
	$(V)go build  -ldflags "-s -w -X main.version=${RELEASE} -X main.buildTime=${BUILD_TIME}" -o $(BOT_OUT) $(BOT_MAIN)
	@echo DONE

#### store API service build for linux
.PHONY: bot-linux
bot-linux: export GOOS := linux
bot-linux: export GOARCH := amd64
bot-linux: bot


#### GOPRIVATE setup https://gist.github.com/MicahParks/1ba2b19c39d1e5fccc3e892837b10e21
GOPRIVATE="github.com/*"
.PHONY: tidy
tidy:
	$(V)GOPRIVATE=$(GOPRIVATE) go mod tidy -v

.PHONY: lint
lint:
	$(V)golangci-lint run --config scripts/.golangci.yml
