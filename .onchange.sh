#!/bin/bash -e
path=$1
dir=$(dirname "$path")
filename=$(basename "$path")
extension="${filename##*.}"
nameonly="${filename%.*}"

case $extension in
    go)
        goimports -w $path
        ;;
esac

go build ./...
go test -coverprofile /tmp/c.out ./...
uncover -min 90 /tmp/c.out
killall tidio
