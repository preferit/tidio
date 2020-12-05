#!/bin/bash -e

case $1 in
    *)
	sudo systemctl stop tidio
	go install ./cmd/...
	sudo systemctl start tidio
	;;
esac

