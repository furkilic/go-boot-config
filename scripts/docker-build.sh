#!/usr/bin/env bash

DIR=`dirname ${0}`
. $DIR/common.sh

export GOOS=linux
$DIR/build.sh

docker build -t go-boot-config:0.0.1 "$BASE_DIR" -f "$BASE_DIR/build/package/Dockerfile"