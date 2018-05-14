#!/usr/bin/env bash

set -x
export GOPATH=$(pwd $(dirname $0))

mkdir -p $GOPATH/src/golang.org/x && cd $GOPATH/src/golang.org/x && git clone https://github.com/golang/image.git

cd ${GOPATH}
go get github.com/afocus/captcha
go get github.com/mattn/go-sqlite3
go get gopkg.in/yaml.v2