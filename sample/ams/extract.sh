#!/bin/bash

echo "#################################################################"
echo "# Извлечь сущность в формате .pem                               #"
echo "#################################################################"

BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")
source "${BIN}/util.sh"

ALIAS=$1
FILE=$2

[ -z "$ALIAS" ] && echo "Не передан alias: ./extract.sh [alias]" && exit 1
[ "$FILE" ] && [ -f "$FILE" ] && echo "Файл ужу существует" && exit 1

runmqakm -cert -extract -db "$_STORE_" -pw "$_PASS_" \
  -label "$ALIAS" \
  -target "$FILE"

echo "$FILE"
