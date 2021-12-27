#!/bin/bash

BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")
source "${BIN:?}/util.sh"

ALIAS=$1
[ -z "$ALIAS" ] && echo "Не передан alias сущности: ../init.sh alias" && exit 1

echo \
  "cms.keystore = путь к директории хранилища + имя файла без расширения
cms.certificate = ${_ALIAS_CERT_}
" >"$_CONF_"

# Create the keystore for CA (/home/mqm/CA)
#######################################
runmqakm -keydb -create -db "$_STORE_" -pw "$_PASS_" \
  -type pkcs12 \
  -expire 1000 \
  -stash

# Create Self Signed CA Certificate
#############################
runmqakm -cert -create -db "$_STORE_" -pw "$_PASS_" \
  -label "$_ALIAS_CERT_" \
  -dn "$_DNAME_" \
  -default_cert yes \
  -sigalg sha256


#runmqakm -cert -extract -db "$_STORE_" -pw "$_PASS_" \
#  -label "$_ALIAS_CERT_" \
#  -target "$_CERT_"


#runmqakm -cert -add -db /home/bob/.mqs/bobkey.kdb -pw passw0rd -label Alice_Cert -file alice_public.arm



#Create the keystore for CA (/home/mqm/CA)
########################################
#1. Sudo runmqckm -keydb -create -db SSL_CA -pw Passw0rd -type cms -expire 365 -
#stash
#
#Create Self Signed CA Certificate
##############################
#2. Sudo runmqckm -cert -create -db SSL_CA.kdb -pw Passw0rd -label ssl_ca -dn
#“CN=SSL CA,O=IBM,C=China” -expire 365
#List Certificate under CA’s Keystore
#################################
#3. Sudo runmqckm -cert -list -db SSL_CA.kdb -pw Passw0rd
#Extract Public (Self Signed) CA Certificate
#######################################
#4. Sudo runmqckm -cert -extract -db SSL_CA.kdb -pw Passw0rd -label ssl_ca -target
#ssl_ca.cer -format ascii
