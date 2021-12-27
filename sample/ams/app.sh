#!/bin/bash

# env файл конфигурации
CFG=$1

# Применить файл конфигурации, если предоставлен
if [ "$CFG" ]; then
  [ ! -f "$CFG" ] && echo "Error: файл не существует" && exit 1




fi


BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")
source "${BIN}/util.sh"


menu() {\
  echo
  echo "Выбор действия: "
  echo ">> 1) Показать значения переменных окружения"
  echo ">> 2) Создать хранилище"
  echo ">> 3) Генерировать пару ключей"
  echo ">> 1"
}


# Установить значения
_DIR_STORE_




LEN=1
MSG_SELECT=""

while true; do
    menu

    read -rn $LEN -p "$MSG_SELECT" ANSWER

    echo ">>>>>>>>>> $ANSWER"
done
