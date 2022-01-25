#!/bin/bash

# Выпуск сертификата

BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")
CSR=
CERT=
PREV=

for it in "$@"; do
  if [ "$PREV" ]; then
    case $PREV in
    "-csr") CSR=$it ;;
    "-cert") CERT=$it ;;
    *) echo "WARN: не распознанные аргументы: $PREV $it" ;;
    esac
    PREV=
    continue
  fi

  PREV=$it
done

unset PREV it

[ ! -f "$CSR" ] &&
  echo "Файл запроса на выпуск сертификата не существует '$CSR'" &&
  exit 1

[ -f "$CERT" ] &&
  echo "Файл сертификата существует '$CERT'" &&
  exit 1

# подпись удостоверяющим центром
openssl x509 -req -days 365 \
  -in "$CSR" \
  -CA "${BIN}/server/cacert.crt" \
  -CAkey "${BIN}/server/ca.key" \
  -CAcreateserial \
  -out "$CERT"
