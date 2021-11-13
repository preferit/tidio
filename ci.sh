#!/bin/bash -e

# This script is meant to provide an easy way to build and test this
# project either from CI pipeline or on a developer computer.

dist=/tmp/tidio

case $1 in
    s|setup)
        go install github.com/gregoryv/stamp/cmd/stamp
	go install github.com/gregoryv/uncover/cmd/uncover
	;;
    b|build)
	mkdir -p $dist
        go generate ./...
        go build -o $dist/tidio ./cmd/tidio
        cp -r ./systemd.service nginx.conf $dist
        cp install.sh $dist
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
    c|clean)
	rm -rf $dist
	;;
    *)
	echo "Usage: $0 [s]etup|[b]uild|[t]est|[u]ncover|[i]nstall|[c]lean"
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
