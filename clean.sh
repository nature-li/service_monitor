#!/usr/bin/env bash

set -x

# set GOPATH
export GOPATH=$(pwd $(dirname $0))

# build Monitor
go clean Monitor
go clean -i Monitor