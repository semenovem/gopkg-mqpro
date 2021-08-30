#!/bin/bash

BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")

set -o allexport
source "${BIN}/project.properties"
set +o allexport

if [ -z "$(docker images --filter=reference="$DOCKER_IMG" -q)" ]; then
  echo
  echo "###########################################################"
  echo "# Сборка образа для разработки с ibmmq                    #"
  echo "###########################################################"

  docker build -f Dockerfile -t "$DOCKER_IMG" ../
fi

echo
echo "###########################################################"
echo "# Старт ibmmq браузера                                    #"
echo "###########################################################"

[ -z "$(docker network ls -f name="$NETWORK" -q)" ] &&
  docker network create --driver overlay --attachable "$NETWORK"

docker run -it --rm \
  --hostname=b \
  --name=browser \
  --network="$NETWORK" \
  -v "${BIN}/connect.env:/app/connect.env:ro" \
  "$DOCKER_IMG" sh -c \
  "echo 'Запустите в отдельном терминале curl' && bash"
