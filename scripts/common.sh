#!/usr/bin/env bash

cygwin=false;
darwin=false;
mingw=false
extension=""

case "`uname`" in
  CYGWIN*) cygwin=true;;
  MINGW*) mingw=true;;
  Darwin*) darwin=true;;
esac


SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
BASE_DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

if $cygwin; then
  BASE_DIR=`cygpath --path --windows "$BASE_DIR"`
fi

if [[ $cygwin == true ]] || [[ $mingw == true ]]; then
  if [ "$GOOS" == "" ]; then
    extension=".exe"
  fi
fi

if [ "$GOOS" == "windows" ] ; then
  extension=".exe"
fi