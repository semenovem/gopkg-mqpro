#!/bin/sh

c () {
  curl client/"$*"
}

c2 () {
  curl client2/"$*"
}
