#!/usr/bin/env bash

set -x

# set GOPATH
export GOPATH=$(pwd $(dirname $0))

# build Monitor
go build -ldflags "-X mt/mtlog.CodeRoot=${GOPATH}" monitor
go install monitor

go build -ldflags "-X mt/mtlog.CodeRoot=${GOPATH}" agent
go install agent
