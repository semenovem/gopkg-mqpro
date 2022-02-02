#!/bin/bash

c () {
  curl client1/"$*"
}

c2 () {
  curl client2/"$*"
}


client () {
  curl client1/"$*"
}

client2 () {
  curl client2/"$*"
}

export -f c c2 client client2
export PS1="> "

echo "-------------------------------------------"
echo "example: "
echo "c get"
echo "c2 put"
echo "client get"
echo "client2 put"

