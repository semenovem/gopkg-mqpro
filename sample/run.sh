#!/bin/bash

BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")

#go get
#
#export GOROOT=/usr/local/go-1.16.4
#
#while true; do
#
#  /usr/local/go-1.16.4/bin/go run *.go
#
#  echo
#  echo "###########################################################"
#  echo "# Старт приложения в DEV MODE                             #"
#  echo "# [ctrl + c] - два раза подряд для выхода                 #"
#  echo "###########################################################"
#  sleep 1
#
#done



if [ "$1" = "run2" ]; then
  echo
  echo "###########################################################"
  echo "# Старт в DEV MODE                                        #"
  echo "# [ctrl + c] - два раза подряд для выхода                 #"
  echo "###########################################################"

  set -o allexport
  source "${BIN}/run.env"
  set +o allexport

  /usr/local/go-1.16.4/bin/go run *.go

  sleep 1
fi



if [ -z "$1" ]; then
  go get

  while true; do
    bash "$0" "run2"
  done
fi
