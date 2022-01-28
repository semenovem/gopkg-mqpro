#!/bin/bash

BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")
app=${BIN:?}/crypto-app.sh

CA_CERT="${BIN:?}/server/cacert.crt"
CA_KEY="${BIN:?}/server/ca.key"

CFG1_KEYSTORE_DIR_NAME="/app/keystore1"
CFG1_KEYSTORE_NAME=mq-ams
CFG1_KEYSTORE_PASSWORD=passw0rd1
CFG1_KEYSTORE_TYPE=cms
CFG1_USER_CERT_LABEL=Alice
CFG1_USER_CERT_DNAME=cn=alice,O=IBM,c=GB
CFG1_USER_CERT_REQ="${CFG1_KEYSTORE_DIR_NAME}/alice-cert-req.pem"
CFG1_USER_CERT="${CFG1_KEYSTORE_DIR_NAME}/alice-cert.pem"
CFG1_CONF_KEYSTORE_PATH=/mqs

CFG2_KEYSTORE_DIR_NAME="/app/keystore2"
CFG2_KEYSTORE_NAME=mq-ams
CFG2_KEYSTORE_PASSWORD=passw0rd2
CFG2_KEYSTORE_TYPE=cms
CFG2_USER_CERT_LABEL=Bob
CFG2_USER_CERT_DNAME=cn=bob,O=IBM,c=GB
CFG2_USER_CERT_REQ="${CFG2_KEYSTORE_DIR_NAME}/bob-cert-req.pem"
CFG2_USER_CERT="${CFG2_KEYSTORE_DIR_NAME}/bob-cert.pem"
CFG2_CONF_KEYSTORE_PATH=/mqs

funcClient1() {
  export MQM_CRYPTO_KEYSTORE_DIR_NAME=$CFG1_KEYSTORE_DIR_NAME
  export MQM_CRYPTO_KEYSTORE_NAME=$CFG1_KEYSTORE_NAME
  export MQM_CRYPTO_KEYSTORE_PASSWORD=$CFG1_KEYSTORE_PASSWORD
  export MQM_CRYPTO_KEYSTORE_TYPE=$CFG1_KEYSTORE_TYPE
  export MQM_CRYPTO_USER_CERT_LABEL=$CFG1_USER_CERT_LABEL
  export MQM_CRYPTO_USER_CERT_DNAME=$CFG1_USER_CERT_DNAME
  export MQM_CRYPTO_USER_CERT_REQ=$CFG1_USER_CERT_REQ
  export MQM_CRYPTO_USER_CERT=$CFG1_USER_CERT
  export MQM_CRYPTO_USER_CERT=$CFG1_CONF_KEYSTORE_PATH
  export MQM_CRYPTO_CONF_KEYSTORE_PATH=$CFG1_CONF_KEYSTORE_PATH
}

funcClient2() {
  export MQM_CRYPTO_KEYSTORE_DIR_NAME=$CFG2_KEYSTORE_DIR_NAME
  export MQM_CRYPTO_KEYSTORE_NAME=$CFG2_KEYSTORE_NAME
  export MQM_CRYPTO_KEYSTORE_PASSWORD=$CFG2_KEYSTORE_PASSWORD
  export MQM_CRYPTO_KEYSTORE_TYPE=$CFG2_KEYSTORE_TYPE
  export MQM_CRYPTO_USER_CERT_LABEL=$CFG2_USER_CERT_LABEL
  export MQM_CRYPTO_USER_CERT_DNAME=$CFG2_USER_CERT_DNAME
  export MQM_CRYPTO_USER_CERT_REQ=$CFG2_USER_CERT_REQ
  export MQM_CRYPTO_USER_CERT=$CFG2_USER_CERT
  export MQM_CRYPTO_CONF_KEYSTORE_PATH=$CFG2_CONF_KEYSTORE_PATH
}

# preparing
# ----------------------------------------------------------
rm -rf $CFG1_KEYSTORE_DIR_NAME
rm -rf $CFG2_KEYSTORE_DIR_NAME
rm -rf "${BIN:?}/server"

mkdir "${BIN:?}/server"

# processing
# ----------------------------------------------------------
openssl req -newkey rsa:2048 -nodes \
  -keyout "$CA_KEY" \
  -x509 -days 365 \
  -out "$CA_CERT" \
  -subj "/C=RU/ST=MO/L=Moscow/O=TEST/OU=Finance/CN='mqm-sample'/emailAddress=email@test.ru"

funcClient1
$app -y -cmd init
$app -y -cmd self
$app -y -cmd ls
$app -y -cmd extract -cmd-label "$CFG1_USER_CERT_LABEL" -cmd-file "$CFG1_USER_CERT"

funcClient2
$app -y -cmd init
$app -y -cmd self
$app -y -cmd ls
$app -y -cmd extract -cmd-label "$CFG2_USER_CERT_LABEL" -cmd-file "$CFG2_USER_CERT"

# Добавление доверенных сертификатов
# ----------------------------------------------------------
$app -y -cmd ca -cmd-label "$CFG1_USER_CERT_LABEL" -cmd-file "$CFG1_USER_CERT"
$app -y -cmd ca -cmd-label "ca-cert" -cmd-file "$CA_CERT"

funcClient1
$app -y -cmd ca -cmd-label "$CFG2_USER_CERT_LABEL" -cmd-file "$CFG2_USER_CERT"
$app -y -cmd ca -cmd-label "ca-cert" -cmd-file "$CA_CERT"
