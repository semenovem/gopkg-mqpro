#!/bin/bash

BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")
TMP=$(mktemp) || exit 1
NUM=$(grep -n '#!/bin/bash' "$0" | tail -1 | grep -Eo '[0-9]+')
CFG="${BIN:?}/../cfg/$1"

[ -z "$NUM" ] && echo "Ошибка при получении номера строки" && exit 1
[ ! -f "$TMP" ] && echo "Не создан временный файл" && exit 1
[ ! -f "$CFG" ] && echo "нет файла конфигурации: '$CFG'" && exit 1

# The line number at which the contents of the temporary file begin
sed -n "${NUM},100p" "$0" > "$TMP" || exit 1

go get || exit 1

while true; do
  bash "$TMP" "$CFG"
  sleep 2
done

exit 0

# Код временного файла для запуска приложения
# ------------------------------------
#!/bin/bash
BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")
echo
echo "###########################################################"
echo "# Старт в DEV MODE                                        #"
echo "# [ctrl + c] - два раза подряд для выхода                 #"
echo "# файл конфигурации = $CFG"
echo "###########################################################"
set -o allexport
source "$1"
set +o allexport
/usr/local/go-1.16.4/bin/go run *.go
