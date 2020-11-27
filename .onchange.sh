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
go build ./...
go test -coverprofile /tmp/c.out $run ./... 
uncover -min 90 /tmp/c.out
go install ./cmd/tidup
go build -o /home/gregory/bin/tidio ./cmd/tidio
killall tidio
echo -e "\033[42m  \e[0m"
