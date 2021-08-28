#!/bin/bash

BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")

set -o allexport
source "${BIN}/project.properties"
set +o allexport

docker run -it --rm \
  --network "$NETWORK" \
  ubuntu:20.04 sh -c \
  "apt update && apt install -y curl && echo '| Для выхода из контейнера ctrl+D' && bash"
