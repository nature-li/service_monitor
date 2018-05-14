#!/usr/bin/env bash

set -x

# set GOPATH
export GOPATH=$(pwd $(dirname $0))

# build http_server
go build -ldflags "-X mtlog.CodeRoot=${GOPATH}" http_server
go install http_server

# collect target files
TARGET_PATH=${GOPATH}/target
TARGET_FILE_SERVER=${TARGET_PATH}/file_server
mkdir -p ${TARGET_PATH}

# bin
mkdir -p ${TARGET_FILE_SERVER}/bin
cp -rf ${GOPATH}/bin/http_server ${TARGET_FILE_SERVER}/bin

# config
mkdir -p ${TARGET_FILE_SERVER}/config
cp -rf ${GOPATH}/config/conf.template.yaml ${TARGET_FILE_SERVER}/config

# templates
mkdir -p ${TARGET_FILE_SERVER}/templates
cp -rf ${GOPATH}/templates/* ${TARGET_FILE_SERVER}/templates

# readme
cp -rf ${GOPATH}/Readme.md ${TARGET_FILE_SERVER}

# mk tar.gz
cd ${TARGET_PATH}
tar -czvf file_server.tar.gz file_server
cd ${GOPATH}

