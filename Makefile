.PHONY: run test deploy help clean

SHELL=/bin/bash -o pipefail
COMMIT_ID=$(shell git rev-parse --short HEAD)
ENVIRONMENT_NAME := $(or ${ENVIRONMENT_NAME},${ENVIRONMENT_NAME},'staging')
AWS_REGION := $(or ${DEFAULT_AWS_REGION},${AWS_REGION},'us-east-1')
PKG_LIST := $(shell go list | grep -v /vendor/)
AWS_ACCOUNT_ID := $(or ${AWS_ACCOUNT_ID}, '381725629183')
PWD:=$(shell pwd)
ENV=local

help: ### Show Help
	@printf "$(OK_COLOR)"
	@echo "=========================================="
	@echo "============  AWAYMAIL API ================"
	@echo "=========================================="
	@printf "$(NO_COLOR)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

lint: ## Lint Golang files
	@golint -set_exit_status ./...

build: ## Build the binary file
	@go build -o am_server

docker-login:
	@bash -c 'aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com'

docker-build: ## Build's Docker Image
	@docker build --build-arg ENV=dev -t cloudsrc/am-api:$(COMMIT_ID) .
	@docker build -t cloudsrc/gaurun-awaymail:$(COMMIT_ID) -f gaurun/Dockerfile gaurun/
	@docker build -t cloudsrc/gaurun-awaymail-dev:$(COMMIT_ID) -f gaurundev/Dockerfile gaurundev/

docker-build-local:
	@docker build --build-arg ENV=local -t cloudsrc/am-api:$(COMMIT_ID) .
	@docker build -t cloudsrc/gaurun-awaymail:$(COMMIT_ID) -f gaurun/Dockerfile gaurun/

docker-build-dev:
	@docker build --build-arg ENV=dev -t cloudsrc/am-api:$(COMMIT_ID) .
	@docker build -t cloudsrc/gaurun-awaymail:$(COMMIT_ID) -f gaurun/Dockerfile gaurun/

docker-build-staging:
	@docker build --build-arg ENV=staging -t cloudsrc/am-api:$(COMMIT_ID) .
	@docker build -t cloudsrc/gaurun-awaymail:$(COMMIT_ID) -f gaurun/Dockerfile gaurun/

docker-build-prod:
	@docker build --build-arg ENV=prod -t cloudsrc/am-api:$(COMMIT_ID) .
	@docker build -t cloudsrc/gaurun-awaymail:$(COMMIT_ID) -f gaurun/Dockerfile gaurun/

docker-push: ## Push to registry
	@docker push cloudsrc/am-api:$(COMMIT_ID)
	@docker push cloudsrc/gaurun-awaymail:$(COMMIT_ID)

docker-build-and-push: docker-login docker-build docker-push 

run: ## Run's Awaymail API in Development Mode
	@realize start

# bring up backend and api container
local-up:
	@docker-compose up -d mongodb rabbitmq redis gaurun
	@sleep 10 # wait for the backend to come up
	@docker run -d -p 8081:8081 --name=am-api --net=am_network cloudsrc/am-api:$(COMMIT_ID)

# remove api and backend containers
local-down:
	@docker container rm -f $$(docker container ls -aq --filter name=am-api)
	@docker-compose down -v --remove-orphans

clean: ## Remove's the docker container to prevent space on system.
	@docker stop `docker ps -a -q` >/dev/null 2>&1 || true
	@docker rm `docker ps -a -q` >/dev/null 2>&1 || true

test: ## Run's the Casbu Test Suite.
	@echo "=================================================================================="
	@echo "Coverage Test"
	@echo "=================================================================================="
	go fmt ./... && go test -coverprofile coverage.cov -cover ./... # use -v for verbose
	@echo "\n"
	@echo "=================================================================================="
	@echo "All Package Coverage"
	@echo "=================================================================================="
	go tool cover -func coverage.cov

test-coverage: clean ## Run tests with coverage
	@GOLANG_ENV=test APP_ROLE=standalone go test -v .../.. -count=1 -coverprofile cover.out -covermode=atomic ${PKG_LIST}
	@cat cover.out >> coverage.txt
	@make clean

deploy: ## Deploy's Awaymail to provided environment.
	@kubectl set image deployment/awaymail awaymail=cloudsrc/am-api:${COMMIT_ID} -n ${ENVIRONMENT_NAME}
	@kubectl set image deployment/gaurun gaurun=cloudsrc/gaurun-awaymail:${COMMIT_ID} -n ${ENVIRONMENT_NAME}
	@kubectl set image deployment/gaurundev gaurundev=cloudsrc/gaurun-awaymail-dev:${COMMIT_ID} -n ${ENVIRONMENT_NAME}