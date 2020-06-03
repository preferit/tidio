#!/bin/bash -e
path=$1
dir=$(dirname "$path")
filename=$(basename "$path")
extension="${filename##*.}"
nameonly="${filename%.*}"

case $extension in
    go)
        goimports -w $path
        gofmt -w $path
        ;;
esac

case $dir in
    website)
	go run ./website -o /tmp/tidiowebsite
	reloadsurf.sh
	;;
esac

go build ./...
go test -coverprofile /tmp/c.out ./...
uncover -min 80 /tmp/c.out
