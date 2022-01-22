DOCKER_IMAGE_NAME := "pkg-mqpro-sample"
VERSION := "1.0"

IMG := "$(DOCKER_IMAGE_NAME):$(VERSION)"
IMG_DEV := "$(IMG)-dev"

PROJ_NAME := "gopkg-mqpro-sample"
NETWORK := "pkg-mqpro-sample"

IMG_CURL := "curl:1"

SHELL = /bin/bash

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
	@$(call checkCreateDockerNet)


define checkCreateDockerNet
		$(eval TTT=$(shell docker network ls -f name=$(NETWORK) -q))
		[ -z "$(TTT)" ] \
		&& docker network create --driver overlay --attachable $(NETWORK) \
		|| echo "docker network $(NETWORK) is exists"
endef


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
 		-f ../docker/curl.dockerfile \
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
	@make CONFIG=client.env NAME_CONTAINER=client HOSTNAME=c KEYSTORE=$(PWD)/ams/keystore1 \
		run_dev_container

## dev2 - контейнер для разработки. Запуск приложения внутри контейнера: `dev/run.sh client2.env`
.PHONY: dev2
dev2:
	@make CONFIG=client2.env NAME_CONTAINER=client2 HOSTNAME=c2 KEYSTORE=$(PWD)/ams/keystore2 \
		run_dev_container


.PHONY: run_dev_container
run_dev_container:
	@docker run --rm -it \
		--name $(NAME_CONTAINER) \
		--hostname $(HOSTNAME) \
		--network $(NETWORK) \
		-w /app/sample \
		-v "${PWD}/..:/app:rw" \
		-v "$(KEYSTORE):/mqs:rw" \
		"$(IMG)" bash
		#dev/run.sh $(CONFIG)


## test
.PHONY: test
test:
	@docker run --rm -it \
		-w /app \
		-v "${PWD}/..:/app:rw" \
		"$(IMG)" bash


## up: launch stack of apps
.PHONY: up
up: net
	@docker stack deploy -c docker-stack.yaml mqpro

## down: launch stack of apps
.PHONY: down
down:
	@docker stack rm mqpro
