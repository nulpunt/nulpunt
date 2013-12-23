#!/bin/sh
export GOPATH=$(pwd)/gopath
echo "Starting npserver compilation"
go build npserver
echo "Compilation done"
./npserver