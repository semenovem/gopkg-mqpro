#!/bin/bash

_RED_='\033[0;31m'
_GREEN_='\033[0;32m'
_YELLOW_='\033[1;33m'
_BLUE_='\033[0;34m'
_PURPLE_='\033[0;35m'
_CYAN_='\033[0;36m'
_LIGHT_GRAY_='\033[0;37m'
_DARK_GRAY_='\033[1;30m'
_LIGHT_RED_='\033[1;31m'
_LIGHT_GREEN_='\033[1;32m'
_LIGHT_BLUE_='\033[1;34m'

_NC_='\033[0m' # No Color

Red_() {
  echo -e "${_RED_}$*${_NC_}"
}

LRed_() {
  echo -e "${_LIGHT_RED_}$*${_NC_}"
}

Gree() {
  echo -e "${_GREEN_}$*${_NC_}"
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
