### imbmq provider sample

## Быстрый старт
подготовка: 

- развернуть docker swarm
- сборка контейнеров
make build-image-curl
make docker


```
cd sample

1. сгенерировать криптоматериалы
make crypto
в файлах: sample/crypto/{keystore1,keystore2}/keystore.conf
исправить строку `cms.keystore = /mq-ams` 
на `cms.keystore = /mqs/mq-ams`


  
2. запустить стек менеджеров ibmmq
make up
  
  
2.1. 
Запустить ams на очередях
make ams
  
  
3. в отдельном терминале - приложение примера использования
make dev
make dev2
  
  
4) в отдельном терминале - контейнер, подключенный к сети приложения для curl запросов
make curl
в контейнере с curl:
```
# Если запущен `make dev`

curl client1/put
curl client1/get
curl client1/browse


# Если запущен `make dev2`

curl client2/put
curl client2/get
curl client2/browse

```
