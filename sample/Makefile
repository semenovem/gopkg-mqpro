DOCKER_IMAGE_NAME := "pkg-mqpro-sample"
VERSION := "1.0"

IMG := "$(DOCKER_IMAGE_NAME):$(VERSION)"
IMGrun_dev_container := "$(IMG)-dev"

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


## curl - launching a container attached to the docker network to execute curl requests
.PHONY: curl
curl:
	@docker run -it --rm \
	-v "${PWD}/cfg/curl-source.sh:/shell.sh:ro" \
	--network $(NETWORK) \
	$(IMG_CURL) sh -c \
     "source /shell.sh && echo 'source /shell.sh' && echo '| Для выхода из контейнера ctrl+D' && sh"


## dev - контейнер для разработки. Запуск приложения внутри контейнера: `dev/run.sh client.env`
.PHONY: dev
dev:
	@make CONFIG=client.env NAME_CONTAINER=client2 HOSTNAME=c run_dev_container

## dev2 - контейнер для разработки. Запуск приложения внутри контейнера: `dev/run.sh client2.env`
.PHONY: dev2
dev2:
	@make CONFIG=client2.env NAME_CONTAINER=client2 HOSTNAME=c run_dev_container


.PHONY: run_dev_container
run_dev_container:
	@docker run --rm -it \
		--name $(NAME_CONTAINER) \
		--hostname $(HOSTNAME) \
		--network $(NETWORK) \
		-w /app/sample \
		-v "${PWD}/..:/app:rw" \
		"$(IMG)" bash dev/run.sh $(CONFIG)


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
