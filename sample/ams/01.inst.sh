#!/bin/bash

echo "#################################################################"
echo "# 01#"
echo "#################################################################"

BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")
source "${BIN:?}/util.sh" || exit 1

runmqakm -keydb -create -db keystore1/key.kdb -type cms -pw passw0rd -stash
