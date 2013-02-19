#!/usr/bin/env bash

echo export GOPATH=$GOPATH:$(dirname "$(readlink -f "$0")")

# vim: set sw=2 sts=2 : #
