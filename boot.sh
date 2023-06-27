#!/usr/bin/env bash

export GOPATH=$(pwd)/../openlib_golang:$(pwd)/../process_public:$(pwd)
export GOOS="linux"

echo $(go version)
echo "GOOS: $GOOS"
echo "GOPATH: $GOPATH"
