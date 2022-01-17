DOCKER_IMAGE_NAME := "pkg-mqpro-sample"
VERSION := "1.0"

IMG := "$(DOCKER_IMAGE_NAME):$(VERSION)"
IMG_DEV := "$(IMG)-dev"

PROJ_NAME := "gopkg-mqpro-sample"
NETWORK := "pkg-mqpro-sample"

IMG_CURL := "curl:1"

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
 		-f ./Dockerfile \
 		-t "$(IMG)" .


## build-image-curl
.PHONY: build-image-curl
build-image-curl:
	@docker build \
 		-f ./Dockerfile-curl.dockerfile \
 		-t "$(IMG_CURL)" .




## crypto - gen crypto
.PHONY: crypto
crypto:
	@bash crypto/run.sh


## ibmmq
.PHONY: ibmmq
ibmmq:
	@docker run -it --rm \
		--name mqpro-sample-ibmmq \
		--network $(NETWORK) \
		--hostname mq0 \
		--env LICENSE=accept \
		--env MQ_QMGR_NAME=QM1 \
		--env MQ_ADMIN_PASSWORD=22 \
		--env MQ_APP_PASSWORD=11 \
		--env MQ_ENABLE_METRICS=true \
		--env MQ_ENABLE_EMBEDDED_WEB_SERVER=true \
		--publish 9443:9443 \
		-v "${PWD}/cfg/ibmmq/30-config.mqsc:/etc/mqm/30-config.mqsc:ro" \
		-v "${PWD}/cfg/ibmmq/script.sh:/app/script.sh:ro" \
    ibmcom/mq:9.2.2.0-r1


## ibmmqtls
.PHONY: ibmmqtls
ibmmqtls:
	@docker run -it --rm \
		--name mqpro-sample-ibmmq \
		--network $(NETWORK) \
		--hostname mq0 \
		--env LICENSE=accept \
		--env MQ_QMGR_NAME=QM1 \
		--env MQ_ADMIN_PASSWORD=22 \
		--env MQ_APP_PASSWORD=11 \
		--env MQ_ENABLE_METRICS=true \
		--env MQ_ENABLE_EMBEDDED_WEB_SERVER=true \
		--publish 9443:9443 \
		-v "${PWD}/crypto/server:/etc_/mqm/pki/keys/mykey:ro" \
		-v "${PWD}/cfg/ibmmq/30-config.mqsc:/etc/mqm/30-config.mqsc:ro" \
		-v "${PWD}/cfg/ibmmq/script.sh:/app/script.sh:ro" \
    ibmcom/mq:9.2.2.0-r1


## ibmmqtls2
.PHONY: ibmmqtls2
ibmmqtls2:
	@docker run -it --rm \
		--name mqpro-sample-ibmmq2 \
		--network $(NETWORK) \
		--hostname mq2 \
		--env LICENSE=accept \
		--env MQ_QMGR_NAME=QM2 \
		--env MQ_ADMIN_PASSWORD=22 \
		--env MQ_APP_PASSWORD=11 \
		--env MQ_ENABLE_METRICS=true \
		--env MQ_ENABLE_EMBEDDED_WEB_SERVER=true \
		--publish 9444:9443 \
		-v "${PWD}/cfg/ibmmq/30-config.mqsc:/etc_/mqm/30-config.mqsc:ro" \
		-v "${PWD}/cfg/ibmmq/script.sh:/app/script.sh:ro" \
    ibmcom/mq:9.2.2.0-r1




## curl - контейнер, присоединенный к сети, для выполнения curl запросов
.PHONY: curl
curl:
	@docker run -it --rm --network $(NETWORK) ubuntu:20.04 sh -c \
     "apt update && apt install -y curl && echo '| Для выхода из контейнера ctrl+D' && bash"

## curl2 - контейнер, присоединенный к сети, для выполнения curl запросов
.PHONY: curl2
curl2:
	@docker run -it --rm \
	-v "${PWD}/dev/shell.sh:/shell.sh:ro" \
	--network $(NETWORK) \
	$(IMG_CURL) sh -c \
     "source /shell.sh && echo 'source /shell.sh' && echo '| Для выхода из контейнера ctrl+D' && sh"

	-#v "${RWD}/dev/shell.sh:/shell.sh:ro" \

## dev - контейнер для разработки. Запуск приложения внутри контейнера: `dev/run.sh`
.PHONY: dev
dev:
	@docker run --rm -it \
		--name client \
		--hostname c \
		--network $(NETWORK) \
		-w /app/sample \
		-v "${PWD}/..:/app:rw" \
		"$(IMG_DEV)" sh -c 'bash dev/run.sh client.env && bash'
# bash dev/run.sh client.env


## dev2 - контейнер для разработки. Запуск приложения внутри контейнера: `dev/run.sh`
.PHONY: dev2
dev2:
	@docker run --rm -it \
		--name client2 \
		--hostname c2 \
		--network $(NETWORK) \
		-w /app/sample \
		-v "${PWD}/..:/app:ro" \
		"$(IMG)" bash
# bash dev/run.sh client.env


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


## server - контейнер для разработки
.PHONY: server
server:
	@docker run --rm -it \
		--name server \
		--hostname s \
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


## up: launch stack of apps
.PHONY: up
up:
	@docker stack deploy -c docker-stack.yaml mqpro

## down: launch stack of apps
.PHONY: down
down:
	@docker stack rm mqpro
