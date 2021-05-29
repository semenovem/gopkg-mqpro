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


## docker
.PHONY: docker
docker:
	@docker build \
		--no-cache \
 		-f ./Dockerfile1-golang \
 		-t "$(IMG)" ..


## run
.PHONY: run
run: docker
	@docker run -it --rm "$(IMG)"



## sample
.PHONY: sample
sample:
	@docker-compose -f docker-compose.yaml -p $(PROJ_NAME) up --build


## net
.PHONY: net
net:
	@docker network create --driver overlay --attachable $(NETWORK)

## install
.PHONY: install
install:
	@docker pull "ibmcom/mq:9.2.2.0-r1"



## dev - контейнер для разработки. Запуск приложения внутри контейнера: `dev/run.sh`
.PHONY: dev
dev:
	@docker run --rm -it \
		--name sample \
		--hostname sample \
		--network $(NETWORK) \
		-w /app/sample \
		-v `pwd`/..:/app:rw \
		"$(IMG)" bash
