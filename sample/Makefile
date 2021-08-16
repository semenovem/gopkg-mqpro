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


## dev - контейнер для разработки. Запуск приложения внутри контейнера: `dev/run.sh`
.PHONY: dev
dev:
	@docker run --rm -it \
		--name sample \
		--hostname sample \
		--network $(NETWORK) \
		--env ENV_LOG_LEVEL=debug \
		--env ENV_API_PORT=80 \
		--env ENV_MQ_0_HOST=mq0 \
		--env ENV_MQ_0_PORT=1414 \
		--env ENV_MQ_0_MGR=QM1 \
		--env ENV_MQ_0_CHANNEL=DEV.APP.SVRCONN \
		--env ENV_MQ_0_PUT_QUEUE=DEV.QUEUE.1 \
		--env ENV_MQ_0_GET_QUEUE=DEV.QUEUE.1 \
		--env ENV_MQ_0_BROWSE_QUEUE=DEV.QUEUE.1 \
		--env ENV_MQ_0_HEADER=RFH2 \
		--env ENV_MQ_0_APP=gopack.mqrpo.sample \
		--env ENV_MQ_0_USER=app \
		--env ENV_MQ_0_PASS=11 \
		--env ENV_MQ_0_TLS=true \
		--env ENV_MQ_0_KEY_REPOSITORY="crypto/client/keys" \
		-w /app/sample \
		-v "${PWD}/..:/app:rw" \
		"$(IMG)" bash run.sh


## curl - контейнер, присоединенный к сети, для выполнения curl запросов
.PHONY: curl
curl:
	@docker run -it --rm --network $(NETWORK) ubuntu:20.04 sh -c \
     "apt update && apt install -y curl && echo '| Для выхода из контейнера ctrl+D' && bash"


## test
.PHONY: test
test:
	@docker run --rm -it \
		-w /app \
		-v "${PWD}/..:/app:rw" \
		"$(IMG)" bash
