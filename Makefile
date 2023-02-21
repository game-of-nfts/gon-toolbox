#!/usr/bin/make -f

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify


build: go.sum
ifeq ($(OS),Windows_NT)
	go build -o build/gontool.exe ./cmd/gontool
else
	go build -o build/gontool ./cmd/gontool
endif

build-linux: go.sum
	GOOS=linux GOARCH=amd64 $(MAKE) build

install: go.sum
	go install  ./cmd/gontool