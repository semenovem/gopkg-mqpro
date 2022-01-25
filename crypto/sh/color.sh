#!/bin/bash

_RED_='\033[0;31m'
_GREEN_='\033[0;32m'
_YELLOW_='\033[1;33m'
_BLUE_='\033[0;34m'
_LIGHT_BLUE_='\033[1;34m'
_PURPLE_='\033[0;35m'
_CYAN_='\033[0;36m'
_LIGHT_GRAY_='\033[0;37m'
_DARK_GRAY_='\033[1;30m'
_LIGHT_RED_='\033[1;31m'
_LIGHT_GREEN_='\033[1;32m'

_NC_='\033[0m' # No Color

_BACKGROUND_BLACK_='\033[40m'
_BACKGROUND_RED_='\033[41m'
_BACKGROUND_GREEN_='\033[42m'
_BACKGROUND_YELLOW_='\033[43m'
_BACKGROUND_DARK_BLUE_='\033[44m'
_BACKGROUND_BLUE_='\033[46m'
_BACKGROUND_PURPLE_='\033[45m'
_BACKGROUND_GRAY_='\033[47m'

Top_() {
  local txt=$* suff
  suff=$(printf '%*s' "$((60 - ${#txt}))" "|")
  echo -e "${_BACKGROUND_DARK_BLUE_}${_YELLOW_}${txt}${suff}${_NC_}"
}

Red_() {
  echo -e "${_RED_}$*${_NC_}"
}

LRed_() {
  echo -e "${_LIGHT_RED_}$*${_NC_}"
}

Gree() {
  echo -e "${_GREEN_}$*${_NC_}"
}

LGre() {
  echo -e "${_LIGHT_GREEN_}$*${_NC_}"
}

Yell() {
  echo -e "${_YELLOW_}$*${_NC_}"
}

Blue() {
  echo -e "${_BLUE_}$*${_NC_}"
}

LBlu() {
  echo -e "${_LIGHT_BLUE_}$*${_NC_}"
}

LGra() {
  echo -e "${_LIGHT_GRAY_}$*${_NC_}"
}

DGra() {
  echo -e "${_DARK_GRAY_}$*${_NC_}"
}

Cyan() {
  echo -e "${_CYAN_}$*${_NC_}"
}

Purp() {
  echo -e "${_PURPLE_}$*${_NC_}"
}

BackgroundRed() {
  echo -e "${_BACKGROUND_RED_}$*${_NC_}"
}

BackgroundGreen() {
  echo -e "${_BACKGROUND_GREEN_}$*${_NC_}"
}

BackgroundYellow() {
  echo -e "${_BACKGROUND_YELLOW_}$*${_NC_}"
}

BackgroundDarkBlue() {
  echo -e "${_BACKGROUND_DARK_BLUE_}$*${_NC_}"
}

BackgroundBlue() {
  echo -e "${_BACKGROUND_BLUE_}$*${_NC_}"
}

BackgroundPurple() {
  echo -e "${_BACKGROUND_PURPLE_}$*${_NC_}"
}

BackgroundGray() {
  echo -e "${_BACKGROUND_GRAY_}$*${_NC_}"
}
