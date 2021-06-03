### imbmq provider sample

## Быстрый старт

```
cd sample
make net

TODO - вынести генерацию криптоматериалов в отдельный контейнер 

// в отдельном терминале - менеджер ibmmq
make ibmmq
bash crypto/gen.sh

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


https://colinpaice.blog/setting-up-tls-for-mq-with-your-own-certificate-authority-using-ikeyman/

https://developer.ibm.com/components/ibm-mq/tutorials/mq-secure-msgs-tls/
https://github.com/ibm-messaging/mq-dev-patterns

set up mutual
https://developer.ibm.com/components/ibm-mq/tutorials/configuring-mutual-tls-authentication-java-messaging-app/

