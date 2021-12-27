#!/bin/bash

BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")
source "${BIN}/util.sh"

echo "#################################################################"
echo "# Извлечь сертификат                                            #"
echo "# alias = ${_ALIAS_CERT_}"
echo "# path  = ${_CERT_}"
echo "#################################################################"

bash extract.sh "$_ALIAS_CERT_" "$_CERT_"
