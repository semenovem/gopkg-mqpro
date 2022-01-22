#!/bin/bash

INFO() {
  LBlu "INFO: $*"
}
WARN() {
  Yell "WARN: $*"
}
ERR() {
  Red_ "ERRO: $*"
}

confirm() {
  local ans msg="Подтвердить ?"
  [ "$1" ] && msg="$1"
  if [ "$_DEV_MODE_" ]; then
    confirmY "$msg"
    return $?
  fi

  while true; do
    read -r -p "$msg [y/N]: " ans
    case "$ans" in
      "y" | "Y" ) return 0 ;;
      "n" | "N" | "") Yell "Отмена действия" && return 1 ;;
    esac
  done
}

confirmY() {
  local ans msg="Подтвердить ?"
  [ "$1" ] && msg="$1"
  while true; do
    read -r -p "$msg [Y/n]: " ans
    case "$ans" in
      "y" | "Y" | "" ) return 0 ;;
      "n" | "N" ) Yell "Отмена действия" && return 1 ;;
    esac
  done
  return 1
}

anyKey() {
  local lab msg="Для продолжения нажмите любою клавишу"
  [ "$1" ] && msg="$1"
  read -rn 1 -p "$msg: " lab
  echo
  return 0
}

selectYesNoQ() {
  local ans msg="Продолжить ?"
  [ "$1" ] && msg="$1"
  while true; do
    read -r -p "$msg [y/n/Q]: " ans
    case "$ans" in
      "y" | "Y" ) return 0 ;;
      "n" | "N" ) return 1 ;;
      "q" | "Q" | "") return 2 ;;
    esac
  done
}


checkKeystore() {
  [ -f "$_KEYSTORE_" ] && return 0
  ERR "файл хранилища не существует: '$_KEYSTORE_'"
  return 1
}

# return 0 - если label есть в списке запросов на выпуск сертификата
# return 1 - отсутствует
checkExistCertReq() {
  local label=$1
  runmqakm -certreq -list -db "$_KEYSTORE_" -stashed | grep "$label" -q
}

listCertsLabels() {
# TODO подготовить список сертификатов
  runmqakm -cert -list -db "$_KEYSTORE_" -stashed
}

# Выбор сертификата из списка хранилища
# Переменная для обмена данными _SHARE_ALIAS_
# TOOD делать выбор через список
pickExistCertificate() {
  # TODO
  # показать список сертификатов
  # выбор сертификата из списка | или ввод его метки

  runmqakm -cert -list -db "$_KEYSTORE_" -stashed || return 1

  enterCertificateLabel true || return 1
  label="$_SHARE_ALIAS_"
}

# Выбор файла через новую оболочку bash
# Переменная для обмена данными _SHARE_PATH_FILE_
pickFile() {
  local tmpFile path newFile ans
  _SHARE_PATH_FILE_=
  tmpFile=$(mktemp)
  [ "$1" ] && newFile=$1

  [ ! -f "$tmpFile" ] && ERR "Системная ошибка. Не создан временный файл" && return 1

  file() {
    local f=$1

    [ -z "$newFile" ] && [ ! -f "$f" ] \
      && echo "ERROR: файл не выбран: '$f'" \
      && return 1

    [[ $f != /* ]] && f="${PWD}/${f}"
    echo "$f" > "$tmpFile"
    exit 0 &> /dev/null
  }

  export -f file
  export tmpFile
  export newFile

  while true; do
    INFO "В новой оболочке наберите file и выберите файл"
    INFO "file /app/folder/file.pem"

    PS1='[example: file /app/folder/file.pem]\$ ' bash || return 1
    path=$(cat "$tmpFile")

    if [ "$newFile" ]; then
      if [ -f "$path" ]; then
          selectYesNoQ "$(WARN "Выбранный файл существует. Перезаписать ?")"
          case "$?" in
            0 ) ;;
            1 ) continue ;;
            2 ) return 1 ;;
          esac
      fi
    else
      [ ! -f "$path" ] \
        && ERR "Ошибка при выборе файла (файл не существует) $path" \
        && return 1
    fi
    break
  done

  _SHARE_PATH_FILE_="$path"
}




# Ввод названия метки сертификата
# Переменная для обмена данными _SHARE_ALIAS_
enterCertificateLabel() {
  local lab ans notConfirm=
  [ "$1" ] && notConfirm=$1

  _SHARE_ALIAS_=
  while [ -z "$_SHARE_ALIAS_" ]; do
    read -r -p "Введите метку сертификата: " lab

    [ -z "$lab" ] && continue
    [ "$notConfirm" ] && _SHARE_ALIAS_="$lab" && return 0

    read -r \
      -p "$(INFO "Метка сертификата: $(Gree "$lab"). Подтвердить ?  [y/n/Q]:")" ans

    case $ans in
      "y" | "Y") _SHARE_ALIAS_="$lab" ;;
      "n" | "N") continue ;;
      "q" | "") Yell "Отмена действия"; return 1 ;;
    esac
  done
}


funcInitKeystore() {
  echo "------------------------------------------------------------"
  Gree "| Инициализировать новое хранилище:"
  DGra "| Директория       | $_KEYSTORE_DIR_"
  DGra "| Файл хранилища   | $_KEYSTORE_"
  DGra "| Тип хранилища    | $_KEYSTORE_TYPE_"
  DGra "| User label       | $_USER_LABEL_"
  echo "------------------------------------------------------------"

  confirm || return 2

  [ -f "$_KEYSTORE_" ] \
    && ERR "файл хранилища уже существует: '$_KEYSTORE_'" \
    && return 1

  mkdir -p "$_KEYSTORE_DIR_" || return 1

  (runmqakm -keydb -create -db "$_KEYSTORE_" -pw "$_KEYSTORE_PASS_" \
    -type "$_KEYSTORE_TYPE_" -stash \
    && INFO "Хранилище создано") \
    || (ERR "Ошибка при создании хранилища" && return 1)

  # Создать файл конфигурации
  echo -e \
    "cms.keystore = /${_KEYSTORE_NAME}\n" \
    "cms.certificate = ${_USER_LABEL_}" > "${_KEYSTORE_DIR_}/keystore.conf"

  INFO "Файл MQ конфигурации: ${_KEYSTORE_DIR_}/keystore.conf"

  return 2
}

funcDestroyKeystore() {
    echo "------------------------------------------------------------"
    Gree "| Удаление хранилища:"
    DGra "| Директория       | $_KEYSTORE_DIR_"
    DGra "| Файл хранилища   | $_KEYSTORE_"
    echo "------------------------------------------------------------"

    LRed_ "Внимание! Будет удалено содержимое директории"

    confirm || return 2

    [ ! -d "$_KEYSTORE_DIR_" ] \
      && ERR "директория хранилища не существует: '$_KEYSTORE_'" \
      && return 1

    checkKeystore || return 1

    (rm -rf "$_KEYSTORE_DIR_" \
      && INFO "Хранилище удалено" \
      && return 2) \
      || (ERR "Хранилище не удалено" && return 1)
}

funcKeystoreList() {
  echo "------------------------------------------------------------"
  Gree "| Содержимое хранилища:"
  DGra "| Файл хранилища   | $_KEYSTORE_"
  echo "------------------------------------------------------------"
  checkKeystore || return 1

  Gree "Сертификаты:"
  runmqakm -cert -list -db "$_KEYSTORE_" -stashed

  Gree "Запросы на выпуск сертификата:"
  runmqakm -certreq -list -db "$_KEYSTORE_" -stashed
}

funcCertificateIssueRequest() {
  echo "------------------------------------------------------------"
  Gree "| Создать запрос на выпуск сертификата:"
  DGra "| Файл хранилища                     | $_KEYSTORE_"
  DGra "| User label                         | $_USER_LABEL_"
  DGra "| DN сертификата                     | $_USER_DNAME_"
  DGra "| Файл запроса на выпуск сертификата | $_USER_CERT_REQ_"
  echo "------------------------------------------------------------"

  confirm || return 2
  checkKeystore || return 1

  [ -f "$_USER_CERT_REQ_" ] \
    && ERR "Файл запроса на выпуск сертификата уже существует: '$_USER_CERT_REQ_'" \
    && return 1

  checkExistCertReq "$_USER_LABEL_" \
    && ERR "Запрос на сертификат '$_USER_CERT_REQ_' уже существует" \
    && return 1

  (runmqakm -certreq -create -db "$_KEYSTORE_" -stashed \
    -label "$_USER_LABEL_" \
    -dn "$_USER_DNAME_" \
    -file "$_USER_CERT_REQ_" \
    && INFO "Запрос на выпуск сертификата создан") \
    || ERR "Запрос на выпуск сертификата не создан"
}

funcAddUserCertificate() {
  local path="$_USER_CERT_"
  echo "------------------------------------------------------------"
  Gree "| Добавить сертификат пользователя:"
  DGra "| Файл хранилища   | $_KEYSTORE_"
  DGra "| User label       | $_USER_LABEL_"
  DGra "| Файл сертификата | $_USER_CERT_"
  echo "------------------------------------------------------------"

  confirm || return 0

  if [ "$path" ]; then
    [ ! -f "$path" ] \
       && WARN "Файл сертификата не существует _USER_CERT_='$_USER_CERT_'" \
       && path=""
  else
    WARN "Файл сертификата не указан _USER_CERT_=''"
  fi

  if [ -z "$path" ]; then
    confirm "хотите указать путь к файлу ?" || return 0
    pickFile || return 1
    path="$_SHARE_PATH_FILE_"
  fi

  (runmqakm -cert -add -db "$_KEYSTORE_" -stashed \
    -label "$_USER_LABEL_" \
    -file "$path" \
  && INFO "Сертификат добавлен") \
  || (ERR "Ошибка добавления сертификата: '$path'" && return 1)
}

funcAddTrustedCertificate() {
  local path label
  echo "------------------------------------------------------------"
  Gree "| Добавить доверенный сертификат (CA или смежных систем):"
  DGra "| Файл хранилища   | $_KEYSTORE_"
  echo "------------------------------------------------------------"

  confirm || return 0

  pickFile || return 1
  path="$_SHARE_PATH_FILE_"

  enterCertificateLabel || return 0
  label="$_SHARE_ALIAS_"

  (runmqakm -cert -add -db "$_KEYSTORE_" -stashed \
    -label "$label" \
    -file "$path" \
  && INFO "Сертификат добавлен") \
  || (ERR "Ошибка добавления сертификата: '$path'" && return 1)
}

funcExtractCertificate() {
  local path label
  echo "------------------------------------------------------------"
  Gree "| Извлечь сертификат в файл .pem:"
  DGra "| Файл хранилища   | $_KEYSTORE_"
  echo "------------------------------------------------------------"

  pickExistCertificate || return 0
  label="$_SHARE_ALIAS_"

  pickFile true || return 1
  path="$_SHARE_PATH_FILE_"

  INFO "label = $label"
  INFO "path  = $path"

  [ -f "$path" ] && (rm -rf "$path" || return 1)

  (runmqakm -cert -extract -db "$_KEYSTORE_" -stashed \
    -label "$label" \
    -target "$path" \
    && INFO "Сертификат экспортирован") \
    || (ERR "Сертификат не экспортирован" && return 1)
}

funcShowCertificate() {
  local label tmpFile
  echo "------------------------------------------------------------"
  Gree "| Показать сертификат: "
  DGra "| Файл хранилища   | $_KEYSTORE_"
  echo "------------------------------------------------------------"

  pickExistCertificate || return 0
  label="$_SHARE_ALIAS_"

  tmpFile=$(mktemp) || return 1
  rm -rf "$tmpFile" || return 1

  runmqakm -cert -extract -db "$_KEYSTORE_" -stashed \
    -label "$label" \
    -target "$tmpFile"

  openssl x509 -noout -text -in "$tmpFile"
}

funcDeleteCertificate() {
  local label
  echo "------------------------------------------------------------"
  Gree "| Удалить сертификат: "
  DGra "| Файл хранилища   | $_KEYSTORE_"
  echo "------------------------------------------------------------"

  pickExistCertificate || return 0
  label="$_SHARE_ALIAS_"

  (runmqakm -cert -delete -db "$_KEYSTORE_" -stashed -label "$label" \
    && INFO "Сертификат удален label=$label") \
    || (ERR "Сертификат не удален label=$label" && return 1)
}

funcCreateSelfSignCertificate() {
  echo "------------------------------------------------------------"
  Gree "| Генерация самоподпсанного сертификата: "
  DGra "| Файл хранилища   | $_KEYSTORE_"
  DGra "| User label       | $_USER_LABEL_"
  DGra "| DN сертификата   | $_USER_DNAME_"
  DGra "| Файл сертификата | $_USER_CERT_"
  echo "------------------------------------------------------------"

  confirm || return 0

  [ -f "$_USER_CERT_" ] \
    && ERR "Файл сертификата существует: '$_USER_CERT_'" \
    && return 1

  (runmqakm -cert -create -db "$_KEYSTORE_" -stashed \
    -label "$_USER_LABEL_" -dn "$_USER_DNAME_" \
    && INFO "Самоподписанный сертификат создан") \
    || (ERR "Самоподписанный сертификат не создан" && return 1)
}




# Произвольная операция с хранилищем
funcArbitraryOperation() {
  echo "------------------------------------------------------------"
  Gree "| Произвольная операция:"
  DGra "| Файл хранилища   | $_KEYSTORE_"
  echo "------------------------------------------------------------"

  [ -f "$_KEYSTORE_" ] || WARN "Хранилища нет"

  echo "ctrl+D | exit - выход"

  _PATH_RUNMQAKM_=$(which runmqakm)
  export _PATH_RUNMQAKM_="$_PATH_RUNMQAKM_"
  export _KEYSTORE_PASS_
  export _KEYSTORE_

  runmqakm() {
    $_PATH_RUNMQAKM_ "$@" -db "$_KEYSTORE_" -stashed
  }
  export -f runmqakm

  PS1='[runmqakm]\$ ' bash
}
