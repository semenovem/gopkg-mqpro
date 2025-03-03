VERSION := "1.0"

IMG := "mqm/sample:$(VERSION)"
IMG_CURL := "mqm/sample/curl:1"
NETWORK := "net-mqm-sample"
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


## build-image
.PHONY: build-image
build-image:
	@docker build -f ./Dockerfile -t "$(IMG)" ./


## build-image-curl
.PHONY: build-image-curl
build-image-curl:
	@docker build -f ../images/curl.dockerfile -t "$(IMG_CURL)" ./


## crypto - gen crypto
.PHONY: crypto
crypto:
	@make -C ../crypto build-image
	@bash crypto/run.sh


## curl - launching a container attached to the docker network to execute curl requests
.PHONY: curl
curl: build-image-curl net
	@docker run -it --rm \
		-v "${PWD}/cfg/curl-source.sh:/curl-source.sh:ro" \
		-v "${PWD}/cfg/smoke-test.sh:/smoke-test.sh:ro" \
		--network $(NETWORK) \
		$(IMG_CURL) \
		bash -c "echo '| Для выхода из контейнера ctrl+D' && source /curl-source.sh && bash"


## dev - контейнер для разработки. Запуск приложения внутри контейнера: `dev/run.sh client1.env`
.PHONY: dev
dev: build-image
	@make CONFIG=client1.env \
		NAME_CONTAINER=client1 \
		HOSTNAME=c \
		KEYSTORE=$(PWD)/crypto/keystore1 \
		run_dev_container


## dev2 - контейнер для разработки. Запуск приложения внутри контейнера: `dev/run.sh client2.env`
.PHONY: dev2
dev2: build-image
	@make CONFIG=client2.env \
		NAME_CONTAINER=client2 \
		HOSTNAME=c2 \
		KEYSTORE=$(PWD)/crypto/keystore2 \
		run_dev_container


.PHONY: run_dev_container
run_dev_container:
	@docker run --rm -it \
		--name $(NAME_CONTAINER) \
		--hostname $(HOSTNAME) \
		--network $(NETWORK) \
		-w /app/sample \
		-v "${PWD}/..:/app:rw" \
		-v "$(KEYSTORE):/mqs:ro" \
		"$(IMG)" bash dev/run.sh $(CONFIG)


# TODO - добавить проверку и запуск: подготовка тестовых криптоматериалов
# загрузку образа ibmmq
## up: launch stack of apps
.PHONY: up
up: net
	@docker stack deploy -c docker-stack.yaml mqm


## down: launch stack of apps
.PHONY: down
down:
	@docker stack rm mqm


## ams - enabled ams
.PHONY: ams
ams:
	@bash dev/ams-run.sh
