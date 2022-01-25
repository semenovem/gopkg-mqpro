#!/bin/sh

c () {
  curl client1/"$*"
}

c2 () {
  curl client2/"$*"
}
