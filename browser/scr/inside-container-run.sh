#!/bin/bash

BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")

while true; do
  bash "${BIN}/inside-container-run2.sh" "$1"
done

