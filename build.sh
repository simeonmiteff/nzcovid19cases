#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

find $DIR/cmd -mindepth 1 -maxdepth 1 -type d | while read D; do
  cd $D
  echo -n "Building app in ${D}... "
  if go build; then
    echo "success."
  else
    echo "failed."
  fi
done