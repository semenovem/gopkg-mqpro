#!/bin/bash

# Запустить скрипты настройки ams на очередях

docker exec -it "$(docker ps -f name=mqm_mq1 -q)" bash /app/cmd.sh
docker exec -it "$(docker ps -f name=mqm_mq2 -q)" bash /app/cmd.sh
