#!/usr/bin/env bash

for p in ${GOPATH//:/ }; do
  [ -d "$p" ] && [[ "$_export" =~ (^|:)$p(:|$) ]] && _export="$_export:$p"
done
echo export GOPATH="$_export":$(dirname "$(readlink -f "$0")")

# vim: set sw=2 sts=2 : #
