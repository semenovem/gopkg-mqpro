#!/bin/bash

BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")
TMP=$(mktemp)

sed -n 19,100p "$0" > "$TMP" # Создает временный файл с содержимым строк 19-

go get

while true; do
  bash "$TMP" "${BIN}/$1"
  sleep 2
done

exit 0

# Код временного файла для запуска приложения
# ------------------------------------

#!/bin/bash

BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")

CFG=$1

echo
echo "###########################################################"
echo "# Старт в DEV MODE                                        #"
echo "# [ctrl + c] - два раза подряд для выхода                 #"
echo "# файл конфигурации = $CFG"
echo "###########################################################"

[ ! -f "$CFG" ] && echo "нет файла конфигурации: '$CFG'" && exit 1

set -o allexport
source "$CFG"
set +o allexport

/usr/local/go-1.16.4/bin/go run *.go
