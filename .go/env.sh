#!/usr/bin/env bash

for p in ${GOPATH//:/ }; do
  [ -d "$p" ] && ! [[ "$_export" =~ (^|:)$p(:|$) ]] && _export="$_export:$p"
done
echo export GOPATH="$_export":$(cd "$(dirname "$0")"; pwd)

# vim: set sw=2 sts=2 : #
