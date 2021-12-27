#!/bin/bash

echo "#################################################################"
echo "# Удалить сущность из хранилища                                 #"
echo "#################################################################"

BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")
source "${BIN}/util.sh"

ALIAS=$1

[ -z "$ALIAS" ] && echo "Не передан alias сущности: ./script alias" && exit 1
[ "$ALIAS" = "$_ALIAS_CERT_" ] && echo "Нельзя удалять приватный ключ" && exit 1

keytool -delete -keystore "$_STORE_" -storepass "$_PASS_" \
    -alias "$ALIAS"
