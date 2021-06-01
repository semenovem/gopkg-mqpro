### imbmq provider sample

## Быстрый старт

```
cd sample
make net

// в отдельном терминале - менеджер ibmmq
make ibmmq

// в отдельном терминале - приложение примера использования
make dev

// в отдельном терминале - контейнер, подключенный к сети приложения для curl запросов
make curl

готово
в контейнере с curl: 

curl sample/get
curl sample/put
curl sample/browse
curl sample/putget
curl sample/sub
curl sample/unsub

```

