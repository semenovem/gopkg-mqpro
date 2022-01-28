#!/bin/bash

BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")
IMG_CRYPTO="mqm/crypto:1"

rm -rf "${BIN:?}/keystore1"
rm -rf "${BIN:?}/keystore2"
rm -rf "${BIN:?}/server"


#docker run -it --rm \
#  -v "${BIN:?}/instruct-gen.sh:/app/instruct-gen.sh:ro" \
#  "$IMG_CRYPTO" \
#  sh -c "bash /app/instruct-gen.sh; bash"
#exit

ID=$(docker run -d \
  -v "${BIN:?}/instruct-gen.sh:/app/instruct-gen.sh:ro" \
  "$IMG_CRYPTO" \
  sh -c "bash /app/instruct-gen.sh; sleep 10") || exit 1

sleep 5

docker cp "${ID:?}:/app/keystore1" "${BIN:?}"
docker cp "${ID:?}:/app/keystore2" "${BIN:?}"
docker cp "${ID:?}:/app/server" "${BIN:?}"

docker stop "$ID"
