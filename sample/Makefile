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
 		-f ./Dockerfile1-golang \
 		-t "$(IMG)" ..


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
		--publish 9443:9443 \
    ibmcom/mq:9.2.2.0-r1


## dev - контейнер для разработки. Запуск приложения внутри контейнера: `dev/run.sh`
.PHONY: dev
dev:
	@docker run --rm -it \
		--name sample \
		--hostname sample \
		--network $(NETWORK) \
		--env ENV_API_PORT=80 \
		--env ENV_MQ_0_HOST=mq0 \
		--env ENV_MQ_0_PORT=1414 \
		--env ENV_MQ_0_MGR=QM1 \
		--env ENV_MQ_0_CHANNEL=DEV.APP.SVRCONN \
		--env ENV_MQ_0_PUT_QUEUE=DEV.QUEUE.1 \
		--env ENV_MQ_0_GET_QUEUE=DEV.QUEUE.1 \
		--env ENV_MQ_0_BROWSE_QUEUE=DEV.QUEUE.1 \
		--env ENV_MQ_0_APP=gopack.mqrpo.sample \
		--env ENV_MQ_0_USER=app \
		--env ENV_MQ_0_PASS=11 \
		-w /app/sample \
		-v `pwd`/..:/app:rw \
		"$(IMG)" go run *.go

