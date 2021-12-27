#!/bin/bash

BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")
source "${BIN:?}/util.sh"


runmqakm -cert -create -db "$_STORE_" -pw "$_PASS_" -label "$_ALIAS_HLDG_" \
  -dn "$_DNAME_" \
  -size 2048 \
  -x509version 3 \
  -expire 365 \
  -fips -sig_alg SHA256WithRSA
