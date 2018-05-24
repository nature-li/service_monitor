#!/usr/bin/env bash

set -x
export GOPATH=$(pwd $(dirname $0))

mkdir -p $GOPATH/src/golang.org/x && cd $GOPATH/src/golang.org/x && git clone https://github.com/golang/image.git

cd ${GOPATH}
go get -u github.com/samuel/go-zookeeper/zk
go get -u -insecure gopkg.in/yaml.v2
go get -u github.com/afocus/captcha

# go get -u golang.org/x/crypto/ssh
git clone https://github.com/golang/crypto.git
cp -rf ${GOPATH}/crypto ${GOPATH}/src/golang.org/x