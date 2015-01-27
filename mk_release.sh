#!/bin/bash

go build -v \
    -ldflags "-X main._VERSION_MINOR $1"\
    -ldflags "-X main._DATE `date -u +%Y-%m-%y.%H:%M:%S`" \
    -ldflags "-X main._COMMIT_ID `git rev-parse --short HEAD`"
