#!/bin/bash

_BIN_=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")

export PATH="${PATH:?}:/opt/mqm/bin"
which runmqakm &> /dev/null || (echo "Нет утилиты 'runmqakm'" && exit 1)

source "${_BIN_:?}/sh/color.sh"
source "${_BIN_:?}/sh/crypto-source.sh"

_DEBUG_=            # Режим отладки
_DEV_MODE_=         # Режим разработки
_SHARE_PATH_FILE_=  # Возвращаемое значение при выборе файла
_SHARE_PATH_DIR_=   # Возвращаемое значение при выборе директории
_SHARE_LABEL_=      # Возвращаемое значение при вводе метки сертификата/запроса
_CONFIRM_YES_=      # Автоматическое согласие при подтверждениях
_CMD_=              # Имя операции для запуска в консольном режиме
_CMD_LABEL_=
_CMD_FILE_=
_CMD_DN_=

# Разбор параметров
prev=
for p in "$@"; do
  if [ "$prev" ]; then
    case $prev in
      "-config") CONFIG_FILE="$p" ;;
      "-cmd") _CMD_="$p" ;;
      "-cmd-label") _CMD_LABEL_="$p" ;;
      "-cmd-file")  _CMD_FILE_="$p" ;;
      "-cmd-dn")    _CMD_DN_="$p" ;;
      * ) WARN "Не известные аргументы: $prev $p" && exit 1
    esac
    prev=
    continue
  fi

  case $p in
    "-debug") _DEBUG_=true ;;
    "-dev-mode") _DEV_MODE_=true ;;
    "-y" | "-yes") _CONFIRM_YES_=true ;;
    *) prev=$p
  esac
done
unset prev p

# Константные значения
_CONST_MENU_MAIN_=main-menu       # Главное меню
_DEFAULT_KEYSTORE_NAME_=keystore  # Дефолтное название хранилища
_DEFAULT_KEYSTORE_TYPE_=cms       # Типа хранилища по умолчанию

# Переменные состояния
_CURRENT_MENU_="$_CONST_MENU_MAIN_" # Текущее открытое меню


# Описание конфигурации
_KEYSTORE_DIR_=       # Путь к директории хранилища
_KEYSTORE_NAME_=      # Имя хранилища
_KEYSTORE_=           # файл хранилища
_KEYSTORE_PASSWORD_=  # Пароль хранилища
_KEYSTORE_TYPE_=      # Тип хранилища
_USER_LABEL_=         # Персональный сертификат пользователя
_USER_DNAME_=         # DN сертификата пользователя
_USER_CERT_REQ_=      # Путь к запросу на выпуск сертификата
_USER_CERT_=          # Сертификат пользователя
_CONF_KEYSTORE_PATH_= # Путь к хранилищу, указываемый в файле конфигурации

# Применить значения из переменных окружения
applyEnvVar() {
  [ -z "$_KEYSTORE_DIR_" ] && [ "$MQM_CRYPTO_KEYSTORE_DIR_NAME" ] \
    && _KEYSTORE_DIR_="$MQM_CRYPTO_KEYSTORE_DIR_NAME"

  [ -z "$_KEYSTORE_NAME_" ] && [ "$MQM_CRYPTO_KEYSTORE_NAME" ] \
    && _KEYSTORE_NAME_="$MQM_CRYPTO_KEYSTORE_NAME"

  [ -z "$_KEYSTORE_PASSWORD_" ] && [ "$MQM_CRYPTO_KEYSTORE_PASSWORD" ] \
    && _KEYSTORE_PASSWORD_="$MQM_CRYPTO_KEYSTORE_PASSWORD"

  [ -z "$_KEYSTORE_TYPE_" ] && [ "$MQM_CRYPTO_KEYSTORE_TYPE" ] \
    && _KEYSTORE_TYPE_="$MQM_CRYPTO_KEYSTORE_TYPE"

  [ -z "$_USER_LABEL_" ] && [ "$MQM_CRYPTO_USER_CERT_LABEL" ] \
    && _USER_LABEL_="$MQM_CRYPTO_USER_CERT_LABEL"

  [ -z "$_USER_DNAME_" ] && [ "$MQM_CRYPTO_USER_CERT_DNAME" ] \
    && _USER_DNAME_="$MQM_CRYPTO_USER_CERT_DNAME"

  [ -z "$_USER_CERT_REQ_" ] && [ "$MQM_CRYPTO_USER_CERT_REQ" ] \
    && _USER_CERT_REQ_="$MQM_CRYPTO_USER_CERT_REQ"

  [ -z "$_USER_CERT_" ] && [ "$MQM_CRYPTO_USER_CERT" ] \
    && _USER_CERT_="$MQM_CRYPTO_USER_CERT"

  [ -z "$_CONF_KEYSTORE_PATH_" ] && [ "$MQM_CRYPTO_CONF_KEYSTORE_PATH" ] \
    && _CONF_KEYSTORE_PATH_="$MQM_CRYPTO_CONF_KEYSTORE_PATH"
}

applyEnvVar

# Если предоставлен файл настроек
if [ "$CONFIG_FILE" ]; then
  if [ -f "$CONFIG_FILE" ]; then
    set -o allexport
    source "$CONFIG_FILE"
    set +o allexport

    applyEnvVar

  else
    ERR "Файл конфигурации не существует. CONFIG_FILE = $CONFIG_FILE"
  fi
fi

unset applyEnvVar

# Установка дефолтных значений
if [ -z "$_KEYSTORE_NAME_" ]; then
  WARN "Не установлено название хранилища [MQM_CRYPTO_KEYSTORE_NAME]. Значение по умолчанию '$_DEFAULT_KEYSTORE_NAME_'"
  _KEYSTORE_NAME_="$_DEFAULT_KEYSTORE_NAME_"
fi

if [ -z "$_KEYSTORE_TYPE_" ]; then
  WARN "Не установлен тип хранилища [MQM_CRYPTO_KEYSTORE_TYPE]. Значение по умолчанию '$_DEFAULT_KEYSTORE_TYPE_'"
  _KEYSTORE_TYPE_="$_DEFAULT_KEYSTORE_TYPE_"
fi

# Файл .sth содержит пароль и работа с хранилищем возможна без его знания
notExistFileSth() {
  [ -f "${_KEYSTORE_DIR_}/${_KEYSTORE_NAME_}.sth" ] && return 1
  return 0
}

# Ввод данных с клавиатуры
enterInitData() {
  local pas1 pas2
  if [ -z "$_KEYSTORE_DIR_" ]; then
    WARN "Не установлен путь к директории хранилища [MQM_CRYPTO_KEYSTORE_DIR_NAME]"
    while true; do
      pickDir || break

      INFO "Путь к директории хранилища: $_SHARE_PATH_DIR_"
      selectYesNoQ
      case $? in
        0) _KEYSTORE_DIR_=$_SHARE_PATH_DIR_; break ;;
        1) continue ;;
        2) break ;;
      esac
    done
  fi

  if [ -z "$_KEYSTORE_PASSWORD_" ] && notExistFileSth; then
    WARN "Не установлен пароль хранилища [MQM_CRYPTO_KEYSTORE_PASSWORD]"
    while true; do
      read -rsp "Пароль: " pas1 && echo
      read -rsp "Повторите пароль: " pas2 && echo

      [ "$pas1" != "$pas2" ] && echo WARN "Пароли не совпадают" && continue
      _KEYSTORE_PASSWORD_="$pas1"
      break
    done
  fi

  if [ -z "$_USER_LABEL_" ]; then
    WARN "Не установлена метка сертификата [MQM_CRYPTO_USER_CERT_LABEL]"
    read -rp "Метка сертификата: " _USER_LABEL_ && echo
  fi

  if [ -z "$_USER_DNAME_" ]; then
    WARN "Не установлено DN (distinguished name) [MQM_CRYPTO_USER_CERT_DNAME]"
    read -rp "Метка сертификата: " _USER_DNAME_ && echo
  fi
}

enterInitData
unset enterInitData

# Установка вычисляемых свойств конфигурации
[[ "$_KEYSTORE_DIR_" != /* ]] && _KEYSTORE_DIR_="${PWD}/${_KEYSTORE_DIR_}"

[ "$_KEYSTORE_DIR_" ] && [ "$_KEYSTORE_NAME_" ] \
  && _KEYSTORE_="${_KEYSTORE_DIR_}/${_KEYSTORE_NAME_}.kdb"

[ -z "$_USER_CERT_REQ_" ] \
  && _USER_CERT_REQ_="${_KEYSTORE_DIR_}/${_USER_LABEL_}-cert-req.pem"

[ -z "$_USER_CERT_" ] && _USER_CERT_="${_KEYSTORE_DIR_}/${_USER_LABEL_}-cert.pem"

_CONF_KEYSTORE_PATH_="${_CONF_KEYSTORE_PATH_}/${_KEYSTORE_NAME_}"

# debug вывод конфигурации
if [ "$_DEBUG_" ]; then
  showCfg() {
    Cyan "$*"
  }
  showCfg "Конфигурация:"
  showCfg "_KEYSTORE_DIR_       = $_KEYSTORE_DIR_"
  showCfg "_KEYSTORE_NAME_      = $_KEYSTORE_NAME_"
  showCfg "_KEYSTORE_           = $_KEYSTORE_"
  showCfg "_KEYSTORE_PASSWORD_  = $_KEYSTORE_PASSWORD_"
  showCfg "_KEYSTORE_TYPE_      = $_KEYSTORE_TYPE_"
  showCfg "_USER_LABEL_         = $_USER_LABEL_"
  showCfg "_USER_DNAME_         = $_USER_DNAME_"
  showCfg "_USER_CERT_REQ_      = $_USER_CERT_REQ_"
  showCfg "_USER_CERT_          = $_USER_CERT_"
  unset showCfg
fi

# Контроль корректности настроек
ERR=

[ -z "$_KEYSTORE_DIR_" ]      && ERR=1 && ERR "не установлено MQM_CRYPTO_KEYSTORE_DIR_NAME"
[ -z "$_KEYSTORE_NAME_" ]     && ERR=1 && ERR "не установлено MQM_CRYPTO_KEYSTORE_NAME"
[ -z "$_KEYSTORE_PASSWORD_" ] \
  && notExistFileSth \
  && ERR=1 && ERR "не установлено MQM_CRYPTO_KEYSTORE_PASSWORD"
[ -z "$_USER_LABEL_" ]        && ERR=1 && ERR "не установлено MQM_CRYPTO_USER_CERT_LABEL"
[ -z "$_USER_DNAME_" ]        && ERR=1 && ERR "не установлено MQM_CRYPTO_USER_CERT_DNAME"

[ "$ERR" ] && exit 100

drawMenuItem() {
  local on=$1 item=$2 num cmd desc

  num=$(echo "$item" | awk '{print $1}')
  cmd=$(echo "$item" | awk '{print $2}')
  desc=$(echo "$item" | awk '{print $3,$4,$5,$6,$7,$8,$9,$10}')

  num=$(printf '%3s' "$num")

  if [ "$on" ]; then
    cmd="[${_LIGHT_GREEN_}${cmd}${_NC_}]"
    cmd=$(printf '%-26s' "$cmd")
  else
    num="${_DARK_GRAY_}${num}"
    cmd="[${cmd}]"
    cmd=$(printf '%-9s' "$cmd")
    desc="${desc}${_NC_}"
  fi

  echo -e "${num} ${cmd} ${desc}"
}

# Главное меню
menu_main () {
  local on no ans
  [ -f "$_KEYSTORE_" ] && on=true
  [ ! -f "$_KEYSTORE_" ] && no=true

  drawMenuItem "$no" "1.  init    Инициализировать хранилище"
  drawMenuItem "$on" "2.  destroy Удалить хранилище"
  drawMenuItem "$on" "3.  ls      Показать содержимое хранилища"
  drawMenuItem "$on" "4.  req     Создать запрос на выпуск сертификата"
  drawMenuItem "$on" "5.  add     Добавить сертификат, выпущенный УЦ по запросу на выпуск сертификата"
  drawMenuItem "$on" "6.  ca      Добавить доверенный сертификат (CA или смежных систем)"
  drawMenuItem "$on" "7.  extract Извлечь сертификат в файл .pem"
  drawMenuItem "$on" "8.  show    Показать данные сертификата"
  drawMenuItem "$on" "9.  rm      Удалить сертификат"

  drawMenuItem "$on" "10. self    Создать самоподписанный сертификат (для тестов)"
  drawMenuItem true "11. exec    Произвольные операции утилитой runmqakm"
  drawMenuItem true "q.  exit    Exit"

  while true; do
    read -r -p "Выбор операции: [номер или $(LGre "command")]: " ans
    menu_main_exec "$ans"
    [ "$?" -eq 10 ] && continue

    while true; do
      echo
      read -rp "[ls,req,add,ca,extract,show,rm] > " ans
      menu_main_exec "$ans"
      [ "$?" -eq 10 ] && return 0
    done
  done
}

# Выполнение действий
menu_main_exec() {
  case $(echo "$1" | awk '{print tolower($0)}') in
    "1"  | "init")    echo; funcInitKeystore ;;
    "2"  | "destroy") echo; funcDestroyKeystore ;;
    "3"  | "ls")      echo; funcKeystoreList ;;
    "4"  | "req")     echo; funcCertificateIssueRequest "$2" ;;
    "5"  | "add")     echo; funcAddUserCertificate ;;
    "6"  | "ca")      echo; funcAddTrustedCertificate ;;
    "7"  | "extract") echo; funcExtractCertificate ;;
    "8"  | "show")    echo; funcShowCertificate ;;
    "9"  | "rm")      echo; funcDeleteCertificate ;;
    "10" | "self")    echo; funcCreateSelfSignCertificate ;;
    "11" | "exec")    echo; funcArbitraryOperation ;;
    "q"  | "exit")    exit 100 ;;
    *) return 10
  esac
}

main() {
  # Консольный режим
  if [ "$_CMD_" ]; then
    [ "$_CMD_LABEL_" ] && _USER_LABEL_=$_CMD_LABEL_
    [ "$_CMD_DN_" ] && _USER_DNAME_=$_CMD_DN_

    menu_main_exec "$_CMD_" "$_CMD_ARGS_"
    return $?
  fi

  # Основной цикл работы
  while true; do
    echo
    echo -e "${_BACKGROUND_PURPLE_}$(printf '%-60s\n' "$_CURRENT_MENU_")${_NC_}"

    case $_CURRENT_MENU_ in
      "$_CONST_MENU_MAIN_") menu_main ;;
      * )  WARN "не валидное значение = '$_CURRENT_MENU_'"  ;;
    esac
  done
}

main
