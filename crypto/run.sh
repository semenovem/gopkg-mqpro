#!/bin/bash

BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")


trap 'echo ""' 2
trap 'echo "Exit"; exit 1' 3

while true
do
  bash "${BIN}/crypto-app.sh" dev-mode

  [ $? -eq "100" ] && exit 0

#  sleep 2
done

exit 0
