#!/bin/bash

go get

export GOROOT=/usr/local/go-1.16.4

while true; do

  /usr/local/go-1.16.4/bin/go run *.go

  echo
  echo "###########################################################"
  echo "# Старт приложения в DEV MODE                             #"
  echo "# [ctrl + c] - два раза подряд для выхода                 #"
  echo "###########################################################"
  sleep 1

done
