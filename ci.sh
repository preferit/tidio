#!/bin/bash -e

# This script is meant to provide an easy way to build and test this
# project either from CI pipeline or on a developer computer.

dist=/tmp/tidio

case $1 in
    s|setup)
	go install github.com/gregoryv/stamp/cmd/stamp@latest
	go install github.com/gregoryv/uncover/cmd/uncover@latest
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
    c|clean)
	rm -rf $dist
	;;
    d|deploy)
	case $TIDIO_HOST in
	    tidio.preferit.se)
		# Guard against releasing untagged
		version=$($dist/tidio --version)
		if [[ "$version" == "unreleased" ]]; then
		    echo "cannot deploy unreleased version to production"
		    exit 0
		fi

		mkdir -p $HOME/.ssh
		go run ./internal/cmd/setupGithubSSH/
		rsync -av --delete-after /tmp/tidio/ $LINODE_USER@$TIDIO_HOST:tidio/
		ssh $LINODE_USER@$TIDIO_HOST 'cd tidio; sudo ./install.sh'
		;;
	    tidio.local)
		# local installation
		pushd $dist
		sudo ./install.sh
		popd
		;;
	    *)
		echo "Missing TIDIO_HOST"
		echo "Use tidio.local or tidio.preferit.se"
		exit 1
		;;
	    esac
	;;
    *)
	echo "Usage: $0 [s]etup|[b]uild|[t]est|[u]ncover|[d]eploy|[c]lean"
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
