#!/bin/bash -e

# This script is meant to provide an easy way to build and test this
# project either from CI pipeline or on a developer computer.

case $1 in
    b|build)
	go build ./...
	;;
    t|test)
	#export group=integration
	go test -coverprofile /tmp/tidio.tprof $run ./...
	;;
    u|uncover)
	uncover -min 77 /tmp/tidio.tprof	
	;;
    i|install)
	# local installation
	sudo systemctl stop tidio
	go install ./cmd/...
	sudo systemctl start tidio
	;;
    *)
	echo "Usage: $0 [b]uild|[t]est|[u]ncover|[i]nstall"
	echo ""
	echo "$0 build test uncover"
	echo "$0 b t u"
	exit 1
	;;
esac


# Run next target if any
shift
[[ -z "$@" ]] && exit 0
$0 $@
