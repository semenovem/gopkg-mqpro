#!/bin/bash

BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")

echo
echo "###########################################################"
echo "# Старт приложения                                        #"
echo "# [ctrl + c] - два раза подряд для выхода                 #"
echo "###########################################################"

set -o allexport
source "${BIN}/../connect.env"
set +o allexport

/app/app-ibmmq-browser

sleep 1
