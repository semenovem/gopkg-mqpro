#!/bin/bash

c put
c2 get

echo "-------------"

c2 put
c get

echo "-------------"

c2 put
c sub
c2 put
c unsub

echo "-------------"

c put
c2 sub
c put
c2 unsub

echo "-------------"

c put
c2 browse
c2 get

echo "-------------"

c2 put
c browse
c get
