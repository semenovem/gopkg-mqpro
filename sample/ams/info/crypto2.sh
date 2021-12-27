#!/bin/bash

# docker run -it --rm -v $PWD:/app -w /app  iata/ibmmq-base:1 bash

# Генерирует криптоматериалы

BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")
export PATH="${PATH}:/opt/mqm/bin"

rm -rf "$BIN"/{mq1,mq2}
mkdir "$BIN"/{mq1,mq2}

#########################################################
# Настройки                                             #
#########################################################
PSW='121212'

MQ1_KDB="${BIN}/mq1/mq1-keys.kdb"
MQ1_DN="cn=hldgappdev201lv.inet.vtb_1,O=VTB,c=RU"
MQ1_LAB="label_mq1"
MQ1_CONF="${BIN}/mq1/keystore.conf"
MQ1_CERT="${BIN}/mq1/cert-mq1.pem"

MQ2_KDB="${BIN}/mq2/mq2-keys.kdb"
MQ2_DN="cn=hldgappdev201lv.inet.vtb_2,O=VTB,c=RU"
MQ2_LAB="label_mq2"
MQ2_CONF="${BIN}/mq2/keystore.conf"
MQ2_CERT="${BIN}/mq2/cert-mq2.pem"

#########################################################
# mq1                                                   #
#########################################################
runmqakm -keydb -create -db "$MQ1_KDB" -pw "$PSW" -stash
runmqakm -cert -create -db "$MQ1_KDB" -pw "$PSW" \
  -label "$MQ1_LAB" -dn "$MQ1_DN" -default_cert yes
echo "cms.keystore = ${MQ1_KDB%.*}
cms.certificate = ${MQ1_LAB}
" >"$MQ1_CONF"

#########################################################
# mq2                                                   #
#########################################################
runmqakm -keydb -create -db "$MQ2_KDB" -pw "$PSW" -stash
runmqakm -cert -create -db "$MQ2_KDB" -pw "$PSW" \
  -label "$MQ2_LAB" -dn "$MQ2_DN" -default_cert yes
echo "cms.keystore = ${MQ2_KDB%.*}
cms.certificate = ${MQ2_LAB}
" >"$MQ2_CONF"

#########################################################
# обмен сертификатами между хранилищами                 #
#########################################################
runmqakm -cert -extract -db "$MQ1_KDB" -pw "$PSW" \
  -label "$MQ1_LAB" -target "$MQ1_CERT"

runmqakm -cert -add -db "$MQ2_KDB" -pw "$PSW" \
  -label "$MQ1_LAB" -file "$MQ1_CERT"

runmqakm -cert -extract -db "$MQ2_KDB" -pw "$PSW" \
  -label "$MQ2_LAB" -target "$MQ2_CERT"

runmqakm -cert -add -db "$MQ1_KDB" -pw "$PSW" \
  -label "$MQ2_LAB" -file "$MQ2_CERT"
