# Copyright (c) 2014, The golo authors
#
GO ?= go
GOIMPORTS ?= goimports
VENDOR := $(CURDIR)/_vendor
BINDIR := ~/bin
GOPATH := $(VENDOR):$(GOPATH)
BUILDNO = $(shell whoami)-$(shell date +%Y%m%d%H%M%S)
PROJ = $(shell basename `pwd`)
COMMITID = $(shell git rev-parse --short HEAD)


all: build

build:
	go build -o ~/bin/golo  -ldflags "-X main.version $(BUILDNO) -X main.commitid $(COMMITID)" ./...
