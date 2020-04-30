#!/usr/bin/env bash

DIR=`dirname ${0}`
. $DIR/common.sh


$BASE_DIR/gow build -o $BASE_DIR/bin/go-boot-config$extension $BASE_DIR/cmd/go-boot-config/main.go