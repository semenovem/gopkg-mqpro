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
  [ "$_CONFIRM_YES_" ] && return 0
  [ "$1" ] && msg="$1"
  if [ "$_DEV_MODE_" ]; then
    confirmY "$msg"
    return $?
  fi

  while true; do
    read -rp "$msg [y/N]: " ans
    case "$ans" in
      "y" | "Y" ) return 0 ;;
      "n" | "N" | "") Yell "Отмена действия" && return 1 ;;
    esac
  done
}

confirmY() {
  local ans msg="Подтвердить ?"
  [ "$_CONFIRM_YES_" ] && return 0
  [ "$1" ] && msg="$1"
  while true; do
    read -rp "$msg [Y/n]: " ans
    case "$ans" in
      "y" | "Y" | "" ) return 0 ;;
      "n" | "N" ) Yell "Отмена действия" && return 1 ;;
    esac
  done
  return 1
}

anyKey() {
  local lab msg
  msg=$(Purp "Для продолжения нажмите любую клавишу")
  [ "$1" ] && msg="$1"
  read -rn 1 -p "$msg: " lab
  echo
  return 0
}

# return 0 - yes
# return 1 - no
# return 2 - Quit
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

# Выбор сертификата из списка хранилища
# Переменная для обмена данными _SHARE_LABEL_
# 0 - все ок
# 1 - ошибка при работе с хранилищем
# 2 - нет сертификатов
# 3 - неизвестная ошибка
pickExistCertificate() {
  local list line arr x mod label dn ans
  _SHARE_LABEL_=

  list=$(runmqakm -cert -list -db "$_KEYSTORE_" -stashed -v | sed 1d)
  [ -z "$list" ] && WARN "Нет сертификатов" && return 2

  Cyan "Список сертификатов в хранилище:"

  x=0
  while read -r line; do
    echo "$line" | grep '\* default' -q && echo "$line" && continue
    ((x++))
    drawItemCert "$line" "$x"
    label=$(echo "$line" | awk '{print $2}')
    arr[$x]="$label"
  done < <(printf '%s\n' "$list")

  echo
  label=
  while [ -z "$label" ]; do
    read -r -p "Выбор сертификата [номер или $(LGre "метка")]: " ans
    [ "${arr[$ans]}" ] && label="${arr[$ans]}" && continue

    for line in ${arr[*]}; do
      [ "$line" == "$ans" ] && label="$line" && break
    done
  done

  [ -z "$label" ] && ERR "Ошибка при выборе сертификата" && return 3

  _SHARE_LABEL_="$label"
}

drawItemCert() {
  local line=$1 num=$2

  mod=$(echo "$line" | awk '{print $1}')
  label=$(echo "$line" | awk '{print $2}')
  dn=$(echo "$line" | awk ' {print $3,$4,$5,$6,$7,$8,$9,$10,$11,$12} ')

  [ "$num" ] && num="${_LIGHT_BLUE_}$(printf '%3s' "${x}.")${_NC_}"
  mod=$(printf '%4s' "[${mod}]")
  label="[${_LIGHT_GREEN_}${label}${_NC_}]"
  dn="${_DARK_GRAY_}${dn}${_NC_}"

  echo -e "${num}${mod}${label} ${dn}"
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
    exit 10 &> /dev/null
  }

  export -f file
  export tmpFile
  export newFile

  while true; do
    INFO "Введите 'file' и выберите файл"

    PS1='[example: file /app/folder/file.pem]\$ ' bash
    case $? in
      10) ;;
      *) return 2 ;;
    esac

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

# Выбор директории
# Переменная для обмена данными _SHARE_PATH_DIR_
# return 0 - ок
# return 1 - ошибка
# return 2 - отмена выбора
pickDir() {
  local tmpFile path ans
  _SHARE_PATH_DIR_=
  tmpFile=$(mktemp)
  [ ! -f "$tmpFile" ] && ERR "Системная ошибка. Не создан временный файл" && return 1

  dir() {
    local d=$1
    [ -z "$d" ] && echo "ERROR: директория не выбрана: '$d'" && return 1
    [ -f "$d" ] && echo "ERROR: выбран файл: '$d'" && return 1
    [[ $d != /* ]] && d="${PWD}/${d}"

    echo "$d" > "$tmpFile"
    exit 0 &> /dev/null
  }

  export -f dir
  export tmpFile

  INFO "Введите 'dir' и выберите/введите существующую или новую директорию"
  PS1='[example: dir /app/folder]\$ ' bash || return 2
  path=$(cat "$tmpFile") || return 1
  [ -z "$path" ] && return 2
  mkdir -p "$path" || return 1

  _SHARE_PATH_DIR_="$path"
}

# Ввод названия метки сертификата
# Переменная для обмена данными _SHARE_LABEL_
# shellcheck disable=SC2120
enterCertificateLabel() {
  local lab ans notConfirm=
  [ "$1" ] && notConfirm=$1

  _SHARE_LABEL_=
  while [ -z "$_SHARE_LABEL_" ]; do
    read -r -p "Введите метку сертификата: " lab

    [ -z "$lab" ] && continue
    [ "$notConfirm" ] && _SHARE_LABEL_="$lab" && return 0

    read -r \
      -p "$(INFO "Метка сертификата: $(Gree "$lab"). Подтвердить ?  [y/n/Q]:")" ans

    case $ans in
      "y" | "Y") _SHARE_LABEL_="$lab" ;;
      "n" | "N") continue ;;
      "q" | "") Yell "Отмена действия"; return 1 ;;
    esac
  done
}

funcInitKeystore() {
  local label=$_USER_LABEL_
  [ "$_CMD_LABEL_" ] && label=$_CMD_LABEL_

  DGra "------------------------------------------------------------"
  Top_ "| Инициализировать новое хранилище:"
  DGra "| Директория       | $_KEYSTORE_DIR_"
  DGra "| Файл хранилища   | $_KEYSTORE_"
  DGra "| Тип хранилища    | $_KEYSTORE_TYPE_"
  DGra "| User label       | $label"
  DGra "------------------------------------------------------------"

  confirm || return 2

  [ -f "$_KEYSTORE_" ] \
    && ERR "файл хранилища уже существует: '$_KEYSTORE_'" \
    && return 1

  mkdir -p "$_KEYSTORE_DIR_" || return 1

  (runmqakm -keydb -create -db "$_KEYSTORE_" -pw "$_KEYSTORE_PASSWORD_" \
    -type "$_KEYSTORE_TYPE_" -stash \
    && INFO "Хранилище создано") \
    || (ERR "Ошибка при создании хранилища" && return 1)

  # Создать файл конфигурации
  if [ "$label" ]; then
    echo "cms.keystore = ${_CONF_KEYSTORE_PATH_}" > "${_KEYSTORE_DIR_}/keystore.conf"
    echo "cms.certificate = ${label}" >> "${_KEYSTORE_DIR_}/keystore.conf"

    INFO "Файл MQ конфигурации: ${_KEYSTORE_DIR_}/keystore.conf"
  fi

  return 0
}

funcDestroyKeystore() {
    DGra "------------------------------------------------------------"
    Top_ "| Удаление хранилища:"
    DGra "| Директория       | $_KEYSTORE_DIR_"
    DGra "| Файл хранилища   | $_KEYSTORE_"
    DGra "------------------------------------------------------------"

    LRed_ "Внимание! Будет удалено содержимое директории"

    confirm || return 2

    [ ! -d "$_KEYSTORE_DIR_" ] \
      && ERR "директория хранилища не существует: '$_KEYSTORE_'" \
      && return 1

    checkKeystore || return 1

    rm -rf "$_KEYSTORE_DIR_" \
      && INFO "Хранилище удалено" \
      && return 2

    ERR "Хранилище не удалено"
    return 1
}

funcKeystoreList() {
  local list
  DGra "------------------------------------------------------------"
  Top_ "| Содержимое хранилища:"
  DGra "| Файл хранилища   | $_KEYSTORE_"
  DGra "------------------------------------------------------------"
  checkKeystore || return 1

  Gree "Сертификаты:"
  list=$(runmqakm -cert -list -db "$_KEYSTORE_" -stashed -v -rfc3339 | sed 1d)
  [ -z "$list" ] && echo "    Нет сертификатов"
  [ "$list" ] && echo -e "$list"

  echo
  Gree "Запросы на выпуск сертификата:"
  list=$(runmqakm -certreq -list -db "$_KEYSTORE_" -stashed -v | sed 1d)
  [ -z "$list" ] && echo "    Нет запросов на выпуск сертификатов"
  [ "$list" ] && echo -e "$list"
}

funcCertificateIssueRequest() {
  local label=$_USER_LABEL_ dn=$_USER_DNAME_ file=$_USER_CERT_REQ_
  [ "$_CMD_FILE_" ] && file=$_CMD_FILE_
  [ "$_CMD_LABEL_" ] && label=$_CMD_LABEL_
  [ "$_CMD_DN_" ] && dn=$_CMD_DN_

  DGra "------------------------------------------------------------"
  Top_ "| Создать запрос на выпуск сертификата:"
  DGra "| Файл хранилища                     | $_KEYSTORE_"
  DGra "| User label                         | $label"
  DGra "| DN сертификата                     | $dn"
  DGra "| Файл запроса на выпуск сертификата | $file"
  DGra "------------------------------------------------------------"

  confirm || return 2
  checkKeystore || return 1

  [ -f "$file" ] \
    && ERR "Файл запроса на выпуск сертификата уже существует: '$file'" \
    && return 1

  checkExistCertReq "$label" \
    && ERR "Запрос на сертификат '$label' уже существует" \
    && return 1

  [[ $file != /* ]] && file="${_KEYSTORE_DIR_}/${file}"

  runmqakm -certreq -create -db "$_KEYSTORE_" -stashed \
    -label "$label" \
    -dn "$dn" \
    -file "$file" \
    && INFO "Запрос на выпуск сертификата создан" \
    && return 0

  ERR "Запрос на выпуск сертификата не создан"
  return 1
}

funcAddUserCertificate() {
  local path="$_USER_CERT_"
  DGra "------------------------------------------------------------"
  Top_ "| Добавить сертификат пользователя:"
  DGra "| Файл хранилища   | $_KEYSTORE_"
  DGra "| User label       | $_USER_LABEL_"
  DGra "| Файл сертификата | $_USER_CERT_"
  DGra "------------------------------------------------------------"

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

  runmqakm -cert -add -db "$_KEYSTORE_" -stashed \
    -label "$_USER_LABEL_" \
    -file "$path" \
  && INFO "Сертификат добавлен" \
  && return 0

  ERR "Ошибка добавления сертификата: '$path'"
  return 1
}

funcAddTrustedCertificate() {
  local path label
  [ "$_CMD_FILE_" ] && path=$_CMD_FILE_
  [ "$_CMD_LABEL_" ] && label=$_CMD_LABEL_

  DGra "------------------------------------------------------------"
  Top_ "| Добавить доверенный сертификат (CA или смежных систем):"
  DGra "| Файл хранилища   | $_KEYSTORE_"
  DGra "------------------------------------------------------------"

  if [ -z "$path" ]; then
    pickFile || return $?
    path="$_SHARE_PATH_FILE_"
  fi

  if [ -z "$label" ]; then
    enterCertificateLabel || return 0
    label="$_SHARE_LABEL_"
  fi

  runmqakm -cert -add -db "$_KEYSTORE_" -stashed \
    -label "$label" \
    -file "$path" \
  && INFO "Сертификат добавлен" \
  && return 0

  ERR "Ошибка добавления сертификата: '$path'"
  return 1
}

funcExtractCertificate() {
  local label=$_CMD_LABEL_ path=$_CMD_FILE_

  DGra "------------------------------------------------------------"
  Top_ "| Извлечь сертификат в файл .pem:"
  DGra "| Файл хранилища   | $_KEYSTORE_"
  DGra "------------------------------------------------------------"

  if [ -z "$label" ]; then
    pickExistCertificate || return 0
    label="$_SHARE_LABEL_"
  fi

  if [ -z "$path" ]; then
    pickFile true || return 1
    path="$_SHARE_PATH_FILE_"
  fi

  INFO "label = $label"
  INFO "path  = $path"

  [ -f "$path" ] && (rm -rf "$path" || return 1)

  runmqakm -cert -extract -db "$_KEYSTORE_" -stashed \
    -label "$label" \
    -target "$path" \
    && INFO "Сертификат экспортирован" \
    && return 0

  ERR "Сертификат не экспортирован"
  return 1
}

funcShowCertificate() {
  local label=$_USER_LABEL_ tmpFile

  DGra "------------------------------------------------------------"
  Top_ "| Показать сертификат: "
  DGra "| Файл хранилища   | $_KEYSTORE_"
  DGra "------------------------------------------------------------"

  if [ -z "$label" ]; then
    pickExistCertificate || return 0
    label="$_SHARE_LABEL_"
  fi

  tmpFile=$(mktemp) || return 1
  rm -rf "$tmpFile" || return 1

  runmqakm -cert -extract -db "$_KEYSTORE_" -stashed \
    -label "$label" \
    -target "$tmpFile" \
  || return 1

  openssl x509 -noout -text -in "$tmpFile"
}

funcDeleteCertificate() {
  local label=$_USER_LABEL_

  DGra "------------------------------------------------------------"
  Top_ "| Удалить сертификат: "
  DGra "| Файл хранилища   | $_KEYSTORE_"
  DGra "------------------------------------------------------------"

  if [ -z "$label" ]; then
    pickExistCertificate || return 0
    label="$_SHARE_LABEL_"
  fi

  runmqakm -cert -delete -db "$_KEYSTORE_" -stashed -label "$label" \
    && INFO "Сертификат удален label=$label" \
    && return 0

  ERR "Сертификат не удален label=$label"
  return 1
}

funcCreateSelfSignCertificate() {
  local label=$_USER_LABEL_ dn=$_USER_DNAME_

  DGra "------------------------------------------------------------"
  Top_ "| Генерация самоподпсанного сертификата: "
  DGra "| Файл хранилища   | $_KEYSTORE_"
  DGra "| User label       | $label"
  DGra "| DN сертификата   | $dn"
  DGra "------------------------------------------------------------"

  confirm || return 0

  runmqakm -cert -create -db "$_KEYSTORE_" -stashed \
    -label "$label" -dn "$dn" \
    && INFO "Самоподписанный сертификат создан" \
    && return 0

  ERR "Самоподписанный сертификат не создан"
  return 1
}

# Произвольная операция с хранилищем
funcArbitraryOperation() {
  local path
  DGra "------------------------------------------------------------"
  Top_ "| Произвольная операция:"
  DGra "| Файл хранилища   | $_KEYSTORE_"
  DGra "------------------------------------------------------------"
  path=$(which runmqakm)
  [ -f "$_KEYSTORE_" ] || WARN "Хранилища нет"
  echo "ctrl+D | exit - выход"

  runmqakm() {
    $_PATH_RUNMQAKM_ "$@" -db "$_KEYSTORE_" -stashed
  }

  export -f runmqakm
  export _KEYSTORE_PASSWORD_
  export _KEYSTORE_

  PS1='[runmqakm]\$ ' _PATH_RUNMQAKM_=$path bash

  unset runmqakm
  return 0
}
