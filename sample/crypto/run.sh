#!/bin/bash

BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")
IMG="pkg-mqpro-sample:1.0"
IMG_CRYPTO="pkg-mqpro-sample/crypto:1.0"

docker build \
  --build-arg BASE_IMAGE="$IMG" \
  -f "${BIN}/Dockerfile" \
  -t "$IMG_CRYPTO" \
  "$BIN"

ID=$(docker run -d "$IMG_CRYPTO" bash /app/gen.sh)

rm -rf "${BIN:?}"/{client,server}

sleep 3

docker cp "${ID:?}:/app/server" "${BIN:?}"
docker cp "${ID:?}:/app/client" "${BIN:?}"

docker stop "$ID"
