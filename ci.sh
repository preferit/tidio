#!/bin/bash -e

# This script is meant to provide an easy way to build and test this
# project either from CI pipeline or on a developer computer.

case $1 in
    build)
	go build ./...
	;;
    test)
	export group=integration
	go test -coverprofile /tmp/tidio.tprof $run ./...
	uncover -min 90 /tmp/tidio.tprof
	;;
    install)
	# local installation
	sudo systemctl stop tidio
	go install ./cmd/...
	sudo systemctl start tidio
	;;
    *)
	echo "Usage: $0 build|test|install"
	exit 1
	;;
esac

