SHELL=/bin/bash
OSNAME=$(shell go env GOOS)
CONFIG_FILE=config.toml
DOCKER_IMG=books_test
DOCKER_TAG=latest
DOCKER_PORT=8080
LOCAL_PORT=8080

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=$(OSNAME) go build -o ./cmd/books/books ./cmd/books

.PHONY: run
run: build
	./cmd/books/books --config=$(CONFIG_FILE)

.PHONY: dockerize
dockerize:
	docker build -t $(DOCKER_IMG):$(DOCKER_TAG) .

.PHONY: run_docker
run_docker: dockerize
	docker run -it \
		-p $(LOCAL_PORT):$(DOCKER_PORT) \
		$(DOCKER_IMG):$(DOCKER_TAG) --config=$(CONFIG_FILE)
