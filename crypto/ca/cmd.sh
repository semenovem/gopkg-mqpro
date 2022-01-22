#!/bin/bash

# Для тестирования
# Удостоверяющий центр bhive

exit

openssl req -x509 -newkey rsa:4096 -days 3650 -nodes \
  -keyout ca-prv.pem \
  -out ca-cert.pem \
  -subj "/C=RU/ST=MO/L=Moscow/O=TEST/OU=Finance/CN='bhive_inet_hldg_20x'/emailAddress=emsemenov@test.ru"

