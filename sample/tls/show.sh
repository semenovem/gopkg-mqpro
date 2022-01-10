#!/bin/bash

echo "#################################################################"
echo "# Показать содержимое сертификата                               #"
echo "#################################################################"

BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")
source "${BIN:?}/util.sh"

ALIAS=$1

[ -z "$ALIAS" ] && echo "Не передан alias сертификата: ./show.sh alias" && exit 1

TMP=$(mktemp) || exit 1

rm -rf "$TMP" || exit 1

runmqakm -cert -extract -db "$_STORE_" -pw "$_PASS_" \
  -label "$ALIAS" \
  -target "$TMP" || exit 1

openssl x509 -noout -text -in "$TMP"
