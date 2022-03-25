.ONESHELL:
# SHELL=C:/Program Files/Git/bin/bash.exe
SHELL=/usr/bin/bash

ifeq ($(OS),Windows_NT)
build: windows
else
build: linux
endif

windows: prepare
	### Set environment variables ###
	export GOOS=windows
	export GOARCH=amd64
	### Build ###
	# go build -ldflags "-s -w" -o ./bin/ip2geo.exe .
	go build -ldflags "-s -w" -o ./ip2geo.exe .

linux: prepare
	### Set environment variables ###
	export GOOS=linux
	export GOARCH=amd64
	### Build ###
	# go build -ldflags "-s -w" -o ./bin/ip2geo .
	go build -ldflags "-s -w" -o ./ip2geo .

prepare:
	### Clear ###
	# rm -rf ./bin

.PHONY: build windows linux
.DEFAULT_GOAL=build
