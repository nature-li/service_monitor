#!/usr/bin/env bash

set -x

# set GOPATH
export GOPATH=$(pwd $(dirname $0))

# build Monitor
go build -ldflags "-X mtlog.CodeRoot=${GOPATH}" Monitor
go install Monitor
