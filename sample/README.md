### imbmq provider sample

## Быстрый старт

```
cd sample

// 1) собрать образ / создать  docker сеть / сгенерировать криптоматериалы
make docker
make net
make crypto

// 2) в отдельном терминале - менеджер ibmmq
make ibmmqtls

// 3) в отдельном терминале - приложение примера использования
make dev

// 4) в отдельном терминале - контейнер, подключенный к сети приложения для curl запросов
make curl

// 5) готово. в контейнере с curl доступны команды: 
curl sample/put
curl sample/get
curl sample/browse
curl sample/putget
curl sample/sub
curl sample/unsub
curl sample/correl
```
