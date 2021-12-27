#!/bin/bash

BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")

DOCKER_IMG="mqpro/ams/gen:2"

if [ -z "$(docker images --filter=reference="$DOCKER_IMG" -q)" ]; then
  echo
  echo "###########################################################"
  echo "# Сборка образа для разработки с ibmmq                    #"
  echo "###########################################################"

  docker build -f Dockerfile -t "$DOCKER_IMG" ./
  [ $? -ne 0 ] && exit 1
fi

echo
echo "###########################################################"
echo "# Старт контейнера для генерации криптоматериалов         #"
echo "###########################################################"

docker run -it --rm -v "${BIN:?}:/app" -w /app "$DOCKER_IMG" bash
