#!/usr/bin/env bash

DIR=`dirname ${0}`
. $DIR/common.sh

docker-compose -f $BASE_DIR/deployments/docker-compose.yml up