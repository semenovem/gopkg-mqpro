#!/bin/bash

# Генерирует криптоматериалы
PSW='&wA*+<_Afh2*4#Z'

BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")
export PATH="${PATH}:/opt/mqm/bin"

rm -rf "${BIN:?}"/{client,server}
mkdir "${BIN:?}"/{client,server}

CACERT="${BIN:?}/server/cacert.crt"
CAKEY="${BIN:?}/server/ca.key"
CLIENT_DB="${BIN:?}/client/keys.kdb"

openssl req -newkey rsa:2048 -nodes -keyout "$CAKEY" -x509 -days 365 -out "$CACERT" \
  -subj "/C=RU/ST=MO/L=Moscow/O=VTB/OU=Finance/CN='mqpro-sample'/emailAddress=email@vtb.ru"

runmqakm -keydb -create -db "$CLIENT_DB" -pw "$PSW" -type pkcs12 -expire 1000 -stash

runmqakm -cert -add -label "QM1.cert" \
  -db "$CLIENT_DB" -stashed \
  -trust enable -file "$CACERT"

sleep 10
