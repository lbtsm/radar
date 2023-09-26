PROJECTNAME=$(shell basename "$(PWD)")
VERSION=-ldflags="-X main.Version=$(shell git describe --tags)"
SOL_DIR=./solidity

CENT_EMITTER_ADDR?=0x1
CENT_CHAIN_ID?=0x1
CENT_TO?=0x1234567890
CENT_TOKEN_ID?=0x5
CENT_METADATA?=0x0

.PHONY: build

build:
	@echo "  >  \033[32mBuilding filter...\033[0m "
	cd cmd && go build -o ../build/filter
