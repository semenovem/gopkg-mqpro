#!/bin/bash


BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")

#CSR=$1
#CERT=$2
#
#[ ! -f "$CSR" ] \
#  && echo "Файл запроса на выпуск сертификата не существует '$CSR'" \
#  && exit 1
#
#[ -f "$CERT" ] \
#  && echo "Файл сертификата существует '$CERT'" \
#  && exit 1

CSR="${BIN}/../keystore1/Alice-cert-req.pem"
CERT="${BIN}/../keystore1/Alice-cert.pem"

ENV_OPENSSL_CFG="${BIN}/ibmmq.conf"


# подпись удостоверяющим центром
openssl x509 -req -days 365 \
  -in "$CSR" \
  -extfile "$ENV_OPENSSL_CFG" \
  -extensions req_ext \
  -CA "${BIN}/ca-cert.pem" \
  -CAkey "${BIN}/ca-prv.pem" \
  -CAcreateserial \
  -out "$CERT"
