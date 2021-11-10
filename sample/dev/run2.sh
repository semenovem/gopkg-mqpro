#!/bin/bash

BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")

ARG1=$1

[ -z "$ARG1" ] && ARG1="client"

echo
echo "###########################################################"
echo "# Старт в DEV MODE                                        #"
echo "# [ctrl + c] - два раза подряд для выхода                 #"
echo "# файл конфигурации = $ARG1"
echo "###########################################################"

CFG="${BIN}/${ARG1}"
[ ! -f "$CFG" ] && echo "нет файла конфигурации: '$CFG'" && exit 1

set -o allexport
source "${BIN}/common-first.env"
source "$CFG"
source "${BIN}/common-last.env"
set +o allexport

/usr/local/go-1.16.4/bin/go run *.go

sleep 1
