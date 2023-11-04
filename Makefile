NAME = elune-backend

VERSION ?= $(shell git symbolic-ref --short -q HEAD)-$(shell git rev-parse --short HEAD)

.DEFAULT_GOAL := help
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

build: clean deps ## Build the project
	go build -ldflags "-s -w" -o bin/$(NAME)

test: deps ## Execute tests
	go test ./...

deps: ## Install dependencies using go get
	go get -d -v -t ./...

clean: ## Remove building artifacts
	rm -rf bin

image: ## Build and push docker image
	docker buildx build --platform linux/arm64,linux/amd64 -t registry.cn-hangzhou.aliyuncs.com/toodo/elune-backend:$(VERSION) . --push

deploy: ## Deploy to k8s
	helm upgrade --install --wait $(NAME) -n toodo ./charts --set image.tag=$(VERSION)