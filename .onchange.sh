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

#run="-run=hacks"
export group=integration
go test -coverprofile /tmp/tidio.tprof $run ./...
uncover -min 90 /tmp/tidio.tprof

