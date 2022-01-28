
SHELL := /bin/bash
IMG := "mqm/crypto:1"


# -----------------
# commands
# -----------------

.PHONY: help
help: Makefile
	@echo "Choose a command in:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'


define checkDockerImage
		$(eval TTT=$(shell docker image ls --filter=reference=$(IMG) -q))
		[ "$(TTT)" ] \
		|| docker build -f ./Dockerfile -t $(IMG) .
endef


## build-image
.PHONY: build-image
build-image:
	@$(call checkDockerImage)


## rm-image
.PHONY: rm-image
rm-image:
	@docker image rm $(IMG)


## dev
.PHONY: dev
dev:
	@docker run -it --rm \
		--name mqm-crypto-dev \
		-w /app \
		-v "$(PWD):/app:rw" \
		$(IMG) sh -c "echo ./crypto-app.sh -config example-config.env; bash"
