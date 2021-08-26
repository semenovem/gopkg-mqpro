#!/bin/bash

BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")

go get

while true; do
  bash "${BIN}/run2.sh" "$1"
done

