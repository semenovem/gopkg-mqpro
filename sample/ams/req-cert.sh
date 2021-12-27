#!/bin/bash

echo "#################################################################"
echo "# Создать запрос на выпуск сертификата                          #"
echo "#################################################################"

BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")
source "${BIN}/util.sh"

[ -f "$_REQ_CERT_" ] && echo "Файл '$_REQ_CERT_' уже существует" && exit 1

keytool -certreq -keystore "$_STORE_" -storepass "$_PASS_" \
    -alias "$_ALIAS_CERT_" \
    -file "$_REQ_CERT_"

[ $? -ne 0 ] && exit 1

echo "Имя файла запроса на выпуск сертификата: $_REQ_CERT_"
