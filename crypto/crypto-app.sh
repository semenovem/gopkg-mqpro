#!/bin/bash

_BIN_=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")

export PATH="${PATH:?}:/opt/mqm/bin"
which runmqakm &> /dev/null || (echo "Нет утилиты 'runmqakm'" && exit 1)

source "${_BIN_:?}/sh/color.sh"
source "${_BIN_:?}/sh/crypto-source.sh"

_DEBUG_=            # Режим отладки
_DEV_MODE_=         # Режим разработки
_SHARE_PATH_FILE_=  # Возвращаемое значение при выборе файла
_SHARE_ALIAS_=      # Возвращаемое значение при вводе метки сертификата/запроса
_CONFIRM_YES_=      # Автоматическое согласие при подтверждениях

# Разбор параметров
for p in "$@"; do
  case $p in
    "debug") _DEBUG_=true ;;
    "dev-mode__") _DEV_MODE_=true ;;
    "-y" | "--yes") _CONFIRM_YES_=true ;;
  esac
done

# Константные значения
_CONST_MENU_MAIN_=main-menu         # Главное меню
_CONST_DEFAULT_KEYSTORE_TYPE_=cms   # Типа хранилища по умолчанию


# Переменные состояния
_CURRENT_MENU_="$_CONST_MENU_MAIN_" # Текущее открытое меню


# Описание конфигурации
_KEYSTORE_DIR_=       # Путь к директории хранилища
_KEYSTORE_NAME=       # Имя хранилища
_KEYSTORE_=           # файл хранилища
_KEYSTORE_PASS_=      # Пароль хранилища
_KEYSTORE_TYPE_=      # Тип хранилища
_USER_LABEL_=         # Персональный сертификат пользователя
_USER_DNAME_=         # DN сертификата пользователя
_USER_CERT_REQ_=      # Путь к запросу на выпуск сертификата
_USER_CERT_=          # Сертификат пользователя

# Подготовка/проверка настроек
# переменные из файла
# переменные окружения

applyEnvVar() {
  if [ -z "$_KEYSTORE_DIR_" ]; then
    if [ "$CFG_KEYSTORE_DIR_NAME" ]; then
      # TODO предусмотреть возможность указания абсолютного пути
      _KEYSTORE_DIR_="${_BIN_}/${CFG_KEYSTORE_DIR_NAME}"
    fi
  fi

  [ -z "$_KEYSTORE_NAME" ] && [ "$CFG_KEYSTORE_NAME" ] \
    && _KEYSTORE_NAME="$CFG_KEYSTORE_NAME"

  [ -z "$_KEYSTORE_PASS_" ] && [ "$CFG_KEYSTORE_PASSWORD" ] \
    && _KEYSTORE_PASS_="$CFG_KEYSTORE_PASSWORD"

  [ -z "$_KEYSTORE_TYPE_" ] && [ "$CFG_KEYSTORE_TYPE" ] \
    && _KEYSTORE_TYPE_="$CFG_KEYSTORE_TYPE"

  [ -z "$_USER_LABEL_" ] && [ "$CFG_USER_CERT_LABEL" ] \
    && _USER_LABEL_="$CFG_USER_CERT_LABEL"

  [ -z "$_USER_DNAME_" ] && [ "$CFG_USER_CERT_DNAME" ] \
    && _USER_DNAME_="$CFG_USER_CERT_DNAME"

  [ -z "$_USER_CERT_REQ_" ] && [ "$CFG_USER_CERT_REQ" ] \
    && _USER_CERT_REQ_="$CFG_USER_CERT_REQ"

  [ -z "$_USER_CERT_" ] && [ "$CFG_USER_CERT" ] \
    && _USER_CERT_="$CFG_USER_CERT"
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

# Установка вычисляемых свойств конфигурации
[ "$_KEYSTORE_DIR_" ] && [ "$_KEYSTORE_NAME" ] \
  && _KEYSTORE_="${_KEYSTORE_DIR_}/${_KEYSTORE_NAME}.kdb"
[ -z "$_USER_CERT_REQ_" ] \
  && _USER_CERT_REQ_="${_KEYSTORE_DIR_}/${_USER_LABEL_}-cert-req.pem"
[ -z "$_USER_CERT_" ] && _USER_CERT_="${_KEYSTORE_DIR_}/${_USER_LABEL_}-cert.pem"
[ -z "$_KEYSTORE_PASS_" ] && _KEYSTORE_PASS_="$_CONST_DEFAULT_KEYSTORE_TYPE_"

showCfg() {
  [ "$_DEBUG_" ] && Cyan "$*"
}

showCfg "Конфигурация:"
showCfg "_KEYSTORE_DIR_   = $_KEYSTORE_DIR_"
showCfg "_KEYSTORE_NAME   = $_KEYSTORE_NAME"
showCfg "_KEYSTORE_       = $_KEYSTORE_"
showCfg "_KEYSTORE_PASS_  = $_KEYSTORE_PASS_"
showCfg "_KEYSTORE_TYPE_  = $_KEYSTORE_TYPE_"
showCfg "_USER_LABEL_     = $_USER_LABEL_"
showCfg "_USER_DNAME_     = $_USER_DNAME_"
showCfg "_USER_CERT_REQ_  = $_USER_CERT_REQ_"
showCfg "_USER_CERT_      = $_USER_CERT_"





# Ввод не необходимых данных

# TODO - пароль и тд


# Контроль корректности настроек
ERR=

[ -z "$_KEYSTORE_DIR_" ]    && ERR=1 && ERR "не установлено CFG_KEYSTORE_DIR_NAME"
[ -z "$_KEYSTORE_NAME" ]    && ERR=1 && ERR "не установлено CFG_KEYSTORE_NAME"
[ -z "$_KEYSTORE_PASS_" ]   && ERR=1 && ERR "не установлено CFG_KEYSTORE_PASSWORD"
[ -z "$_USER_LABEL_" ]      && ERR=1 && ERR "не установлено CFG_USER_CERT_LABEL"
[ -z "$_USER_DNAME_" ]      && ERR=1 && ERR "не установлено CFG_USER_CERT_DNAME"

[ "$ERR" ] && exit 100


# TODO тут место для консольного режима


# Главное меню
menu_main () {
  local cmd=echo ini=echo no="" yes=""
  [ ! -f "$_KEYSTORE_" ] && cmd=DGra && no="__"
  [ -f "$_KEYSTORE_" ] && yes="__" && ini=DGra

  $ini "1.  [init]    Инициализировать хранилище"
  $cmd "2.  [destroy] Удалить хранилище"
  $cmd "3.  [ls]      Показать содержимое хранилища"
  $cmd "4.  [req]     Создать запрос на выпуск сертификата"
  $cmd "5.  [add]     Добавить сертификат, выпущенный УЦ по запросу на выпуск сертификата"
  $cmd "6.  [ca]      Добавить доверенный сертификат (CA или смежных систем)"
  $cmd "7.  [extract] Извлечь сертификат в файл .pem"
  $cmd "8.  [show]    Показать данные сертификата"
  $cmd "9.  [rm]      Удалить сертификат"

  $cmd "10. [self]    Создать самоподписанный сертификат (для тестов)"
  echo "11. [exec]    Произвольные операции утилитой runmqakm"
  echo "q.  [exit]    Exit"

  # Ожидание ввода
  read -r -p "Выбор операции: [command]: " ans
  echo

  menu_main_exec "$ans"
}


# Выполнение действий
menu_main_exec() {
  case $1 in
    "1"  | "init")    funcInitKeystore ;;
    "2"  | "destroy") funcDestroyKeystore ;;
    "3"  | "ls")      funcKeystoreList ;;
    "4"  | "req")     funcCertificateIssueRequest ;;
    "5"  | "add")     funcAddUserCertificate ;;
    "6"  | "ca")      funcAddTrustedCertificate ;;
    "7"  | "extract") funcExtractCertificate ;;
    "8"  | "show")    funcShowCertificate ;;
    "9"  | "rm")      funcDeleteCertificate ;;
    "10" | "self")    funcCreateSelfSignCertificate ;;
    "11" | "exec")    funcArbitraryOperation ;;
    "q"  | "exit")    exit 100 ;;
  esac
}


# Основной цикл работы
while true; do
  echo "************************************************************"
  Gree "$_CURRENT_MENU_"
  echo "************************************************************"

  # Напечатать меню
  case $_CURRENT_MENU_ in
    "$_CONST_MENU_MAIN_") menu_main ;;
    * )  WARN "не валидное значение = '$_CURRENT_MENU_'"  ;;
  esac

  ret=$?

  echo
  case "$ret" in
    0 ) anyKey "$(Yell "Для продолжения нажмите любую клавишу")" ;;
    2 ) sleep 1 ;;
    * ) anyKey ;;
  esac


  break
done
