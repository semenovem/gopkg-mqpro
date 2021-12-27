#!/bin/bash

# docker run -it --rm -v $PWD:/app -w /app  iata/ibmmq-base:1 bash

# Генерирует криптоматериалы
PSW='121212'

BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")
export PATH="${PATH}:/opt/mqm/bin"

rm -rf "$BIN"/{mq1,mq2,ca}
mkdir "$BIN"/{mq1,mq2,ca}

#########################################################
# Настройки                                             #
#########################################################
CA_KEY="${BIN}/ca/ca-prv-key.pem"
CA_CERT="${BIN}/ca/ca-cert.pem"

MQ1_PRV="${BIN}/mq1/mq1-prv-key.pem"
MQ1_CSR="${BIN}/mq1/mq1-csr.pem"
MQ1_KDB="${BIN}/mq1/mq1-keys.kdb"
MQ1_CERT="${BIN}/mq1/mq1-cert.pem"
MQ1_PRV_LABEL="mq1-2021-hldgappdev201lv-prv-key"
MQ1_CERT_LABEL="mq1-2021-hldgappdev201lv-cert"

MQ2_PRV="${BIN}/mq2/mq2-prv-key.pem"
MQ2_CSR="${BIN}/mq2/mq2-csr.pem"
MQ2_KDB="${BIN}/mq2/mq2-keys.kdb"
MQ2_CERT="${BIN}/mq2/mq2-cert.pem"
MQ2_PRV_LABEL="mq2-2021-hldgappdev201lv-prv-key"
MQ2_CERT_LABEL="mq2-2021-hldgappdev201lv-cert"

echo "#########################################################"
echo "# root сертификат                                       #"
echo "#########################################################"
openssl req -x509 -newkey rsa:4096 -days 3650 -nodes \
  -keyout "$CA_KEY" \
  -out "$CA_CERT" \
  -subj "/C=RU/ST=MO/L=Moscow/O=VTB/OU=Finance/CN='bhive_inet_hldg_20x'/emailAddress=emsemenov@vtb.ru"

echo "#########################################################"
echo "# MQ 1                                                  #"
echo "#########################################################"

openssl genrsa -out "$MQ1_PRV" 2048

runmqakm -keydb -create -db "$MQ1_KDB" -pw "$PSW" -type pkcs12 -expire 1000 -stash

runmqakm -cert -add -db "$MQ1_KDB" -pw "$PSW" \
  -file "$MQ1_PRV" -label "$MQ1_PRV_LABEL"

runmqakm -certreq -create -db "$MQ1_KDB" -pw "$PSW" \
  -file "$MQ1_CSR" \
  -label "$MQ1_CERT_LABEL" \
  -dn "C=RU,ST=Moscow,L=Moscow,O=VTB,OU=afsc,CN=hldgappdev201lv.inet.vtb.ru_1" \
  -sigalg sha1

# подпись удостоверяющим центром
openssl x509 -req -days "365" \
  -in "$MQ1_CSR" \
  -extensions req_ext \
  -CA "$CA_CERT" \
  -CAkey "$CA_KEY" \
  -sha1 \
  -CAcreateserial \
  -out "$MQ1_CERT"

# Добавляем доверенный серт
runmqakm -cert -add -db "$MQ1_KDB" -pw "$PSW" \
  -file "$CA_CERT" -label "cacert"

# добавляет выпущенный сертификат в хранилище
runmqakm -cert -receive -db "$MQ1_KDB" -pw "$PSW" \
  -file "$MQ1_CERT"

#runmqakm -cert -list -db "$MQ1_KDB" -pw "$PSW"
#exit

echo "#########################################################"
echo "# MQ 2                                                  #"
echo "#########################################################"

openssl genrsa -out "$MQ2_PRV" 2048

runmqakm -keydb -create -db "$MQ2_KDB" -pw "$PSW" -type pkcs12 -expire 1000 -stash

runmqakm -cert -add -db "$MQ2_KDB" -pw "$PSW" \
  -file "$MQ2_PRV" -label "$MQ2_PRV_LABEL"

runmqakm -certreq -create -db "$MQ2_KDB" -pw "$PSW" \
  -file "$MQ2_CSR" \
  -label "$MQ2_CERT_LABEL" \
  -dn "C=RU,ST=Moscow,L=Moscow,O=VTB,OU=afsc,CN=hldgappdev201lv.inet.vtb.ru_2" \
  -sigalg sha1

# подпись удостоверяющим центром
openssl x509 -req -days "365" \
  -in "$MQ2_CSR" \
  -extensions req_ext \
  -CA "$CA_CERT" \
  -CAkey "$CA_KEY" \
  -sha1 \
  -CAcreateserial \
  -out "$MQ2_CERT"

# Добавляем доверенный серт
runmqakm -cert -add -db "$MQ2_KDB" -pw "$PSW" \
  -file "$CA_CERT" -label "cacert"

# добавляет выпущенный сертификат в хранилище
runmqakm -cert -receive -db "$MQ2_KDB" -pw "$PSW" \
  -file "$MQ2_CERT"

#runmqakm -cert -list -db "$MQ2_KDB" -pw "$PSW"

echo "#########################################################"
echo "#                                                       #"
echo "#########################################################"

runmqakm -cert -add -db "$MQ1_KDB" -pw "$PSW" \
  -file "$MQ2_CERT" -label "$MQ2_CERT_LABEL"

runmqakm -cert -add -db "$MQ2_KDB" -pw "$PSW" \
  -file "$MQ1_CERT" -label "$MQ1_CERT_LABEL"

#runmqakm -cert -list -db "$MQ1_KDB" -pw "$PSW"
#runmqakm -cert -list -db "$MQ2_KDB" -pw "$PSW"

chmod -R 0777 "${BIN}/mq1"
chmod -R 0777 "${BIN}/mq2"
chmod -R 0777 "${BIN}/ca"
