#!/usr/bin/env bash

set -e

CURDIR=`pwd`
OLDGOPATH="$GOPATH"
OLDGOBIN="$GOBIN"
OLDGO111MODULE="$GO111MODULE"
export GOPROXY="https://goproxy.cn"
echo GOPROXY="$GOPROXY"
export GO111MODULE="on"
export GOBIN="$CURDIR/bin/"
echo 'GO111MODULE':$OLDGO111MODULE
echo 'GOPATH:' $GOPATH
echo 'GOBIN:' $GOBIN
echo 'CURDIR:' $CURDIR
go mod vendor
go build -mod=vendor -o go-orm -race -work -v -ldflags "-s" -gcflags "-N -l" cmd/orm/orm.go

if [ ! -d ./bin ]; then
	mkdir bin
fi

if [ -e ./go-orm ]; then
   mv go-orm ./bin/
fi

#if [ -e ./vendor ]; then
#   rm -rf ./vendor
#fi

export GOPATH="$OLDGOPATH"
export GO111MODULE="$OLDGO111MODULE"
export GOBIN="$OLDGOBIN"

echo 'build finished'
