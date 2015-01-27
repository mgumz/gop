#!/bin/bash

D=`date -u +'%Y-%m-%y %H:%M:%S'`
G=`git rev-parse --short HEAD`

go build -v \
    -ldflags "-X main._VERSION_MINOR \"$1\" -X main._DATE \"$D\" -X main._COMMIT_ID \"$G\""
