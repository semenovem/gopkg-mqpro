#!/bin/bash

BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")
source "${BIN:?}/util.sh"

[ -f "$_STORE_" ] && echo "Файл .kdb уже существует" && exit 1

runmqakm -keydb -create -db "$_STORE_" -pw "$_PASS_" -type cms \
  -stash -fips -strong
