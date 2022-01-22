
SHELL := /bin/bash
IMG := "mqpro/crypto:1"


# -----------------
# commands
# -----------------

.PHONY: help
help: Makefile
	@echo "Choose a command in:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'


define checkCreateDockerNet
		$(eval TTT=$(shell docker network ls -f name=$(NETWORK) -q))
		[ -z "$(TTT)" ] \
		&& docker network create --driver overlay --attachable $(NETWORK) \
		|| echo "docker network $(NETWORK) is exists"
endef


## build-image
.PHONY: build-image
build-image:
	@docker build \
 		-f ../docker/crypto.dockerfile \
 		-t "$(IMG)" .

## dev
.PHONY: dev
dev:
	@docker run -it --rm \
		--name mqpro-crypto-dev \
		--entrypoint "" \
		-w /app \
		-v "$(PWD):/app:rw" \
		$(IMG) bash

#bash run.sh


## crypto - gen crypto
.PHONY: crypto
crypto:
	@bash crypto/run.sh
