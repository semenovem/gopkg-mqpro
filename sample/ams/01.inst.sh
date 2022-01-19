#!/bin/bash

echo "#################################################################"
echo "# 01                                                            #"
echo "#################################################################"

BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")
#source "${BIN:?}/util.sh" || exit 1

KEYDB_DIR1="${BIN}/keystore1"
KEYDB1="${KEYDB_DIR1}/mq-ams.kdb"
CFG1="${KEYDB_DIR1}/keystore.conf"
CERT1="${KEYDB_DIR1}/alice_public.arm"

KEYDB_DIR2="${BIN}/keystore2"
KEYDB2="${KEYDB_DIR2}/mq-ams.kdb"
CFG2="${KEYDB_DIR2}/keystore.conf"
CERT2="${KEYDB_DIR2}/bob_public.arm"


rm -rf "$KEYDB_DIR1"
rm -rf "$KEYDB_DIR2"
mkdir -p "$KEYDB_DIR1"
mkdir -p "$KEYDB_DIR2"


runmqakm -keydb -create -db "$KEYDB1" -type cms -pw passw0rd -stash
runmqakm -cert -create -db "$KEYDB1" -pw passw0rd \
  -label Alice_Cert -dn "cn=alice,O=IBM,c=GB" -default_cert yes

echo "cms.keystore = /mqs/mq-ams" > "$CFG1"
echo "cms.certificate = Alice_Cert" >> "$CFG1"


# -----------------
runmqakm -keydb -create -db "$KEYDB2" -type cms -pw passw0rd -stash
runmqakm -cert -create -db "$KEYDB2" -pw passw0rd \
  -label Bob_Cert -dn "cn=bob,O=IBM,c=GB" -default_cert yes

echo "cms.keystore = /mqs/mq-ams" > "$CFG2"
echo "cms.certificate = Bob_Cert" >> "$CFG2"


#---------------------
#---------------------
#---------------------

runmqakm -cert -extract -db "$KEYDB1" -pw passw0rd \
  -label Alice_Cert -target "$CERT1"

runmqakm -cert -add  -db "$KEYDB2" -pw passw0rd \
  -label Alice_Cert -file "$CERT1"

##################################

runmqakm -cert -extract -db "$KEYDB2" -pw passw0rd \
  -label Bob_Cert -target "$CERT2"

runmqakm -cert -add  -db "$KEYDB1" -pw passw0rd \
  -label Bob_Cert -file "$CERT2"
