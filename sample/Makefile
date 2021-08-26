DOCKER_IMAGE_NAME := "pkg-mqpro-sample"
VERSION := "1.0"

IMG := "$(DOCKER_IMAGE_NAME):$(VERSION)"
PROJ_NAME := "gopkg-mqpro-sample"
NETWORK := "pkg-mqpro-sample"


# -----------------
# commands
# -----------------

.PHONY: help
help: Makefile
	@echo "Choose a command in:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'

## net
.PHONY: net
net:
	@docker network create --driver overlay --attachable $(NETWORK)


## install
.PHONY: install
install:
	@docker pull "ibmcom/mq:9.2.2.0-r1"


## docker
.PHONY: docker
docker:
	@docker build \
		--no-cache \
 		-f ./Dockerfile \
 		-t "$(IMG)" .


## crypto - gen crypto
.PHONY: crypto
crypto:
	@bash crypto/run.sh


## ibmmq
.PHONY: ibmmq
ibmmq:
	@docker run -it --rm \
		--name pkg-mqpro-sample-ibmmq \
		--network $(NETWORK) \
		--hostname mq0 \
		--env LICENSE=accept \
		--env MQ_QMGR_NAME=QM1 \
		--env MQ_ADMIN_PASSWORD=22 \
		--env MQ_APP_PASSWORD=11 \
		--env MQ_ENABLE_METRICS=true \
		--env MQ_ENABLE_EMBEDDED_WEB_SERVER=true \
		--publish 9443:9443 \
    ibmcom/mq:9.2.2.0-r1


## ibmmqtls
.PHONY: ibmmqtls
ibmmqtls:
	@docker run -it --rm \
		--name pkg-mqpro-sample-ibmmq \
		--network $(NETWORK) \
		--hostname mq0 \
		--env LICENSE=accept \
		--env MQ_QMGR_NAME=QM1 \
		--env MQ_ADMIN_PASSWORD=22 \
		--env MQ_APP_PASSWORD=11 \
		--env MQ_ENABLE_METRICS=true \
		--env MQ_ENABLE_EMBEDDED_WEB_SERVER=true \
		--publish 9443:9443 \
		-v "${PWD}/crypto/server:/etc/mqm/pki/keys/mykey:ro" \
    ibmcom/mq:9.2.2.0-r1


## curl - контейнер, присоединенный к сети, для выполнения curl запросов
.PHONY: curl
curl:
	@docker run -it --rm --network $(NETWORK) ubuntu:20.04 sh -c \
     "apt update && apt install -y curl && echo '| Для выхода из контейнера ctrl+D' && bash"


## dev - контейнер для разработки. Запуск приложения внутри контейнера: `dev/run.sh`
.PHONY: dev
dev:
	@docker run --rm -it \
		--name client1 \
		--hostname c1 \
		--network $(NETWORK) \
		-w /app/sample \
		-v "${PWD}/..:/app:rw" \
		"$(IMG)" bash dev/run.sh client.env


## dev2 - контейнер для разработки. Запуск приложения внутри контейнера: `dev/run.sh`
.PHONY: dev2
dev2:
	@docker run --rm -it \
		--name client2 \
		--hostname c2 \
		--network $(NETWORK) \
		-w /app/sample \
		-v "${PWD}/..:/app:ro" \
		"$(IMG)" bash dev/run.sh client.env


## dev3 - контейнер для разработки. Запуск приложения внутри контейнера: `dev/run.sh`
.PHONY: dev3
dev3:
	@docker run --rm -it \
		--name client3 \
		--hostname c3 \
		--network $(NETWORK) \
		-w /app/sample \
		-v "${PWD}/..:/app:ro" \
		"$(IMG)" bash dev/run.sh client.env


## push - контейнер для разработки
.PHONY: push
push:
	@docker run --rm -it \
		--name push \
		--hostname push \
		--network $(NETWORK) \
		-w /app/sample \
		-v "${PWD}/..:/app:ro" \
		"$(IMG)" bash dev/run.sh initiator.env


## server - контейнер для разработки
.PHONY: server
server:
	@docker run --rm -it \
		--name push \
		--hostname push \
		--network $(NETWORK) \
		-w /app/sample \
		-v "${PWD}/..:/app:ro" \
		"$(IMG)" bash dev/run.sh server.env


## test
.PHONY: test
test:
	@docker run --rm -it \
		-w /app \
		-v "${PWD}/..:/app:rw" \
		"$(IMG)" bash
