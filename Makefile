SHELL          := /bin/bash
APP_NAME       := aws-ecr-login-secret-updater
TZ             ?= Asia/Tokyo
CGO_ENABLED    ?= 1
AWS_REGION     ?= ap-northeast-1
AWS_ACCOUNT_ID ?= 000000000000
EMAIL          ?= foo@example.com
SECRET         ?= sample-ecr-login-secret
NAMESPACE      ?= default

all: build test lint

build:
	@CGO_ENABLED=${CGO_ENABLED} go build -ldflags="-s -w" -trimpath -tags timetzdata -o ${APP_NAME}

test:
	@go test ./...

lint:
	@go vet ./...
	@golint -set_exit_status ./...

run:
	@TZ=${TZ} \
	AWS_REGION=${AWS_REGION} \
	AWS_ACCOUNT_ID=${AWS_ACCOUNT_ID} \
	EMAIL=${EMAIL} \
	SECRET=${SECRET} \
	NAMESPACE=${NAMESPACE} \
	./${APP_NAME} \
	--kubeconfig=$$HOME/.kube/config

clean:
	@rm -f ${APP_NAME} main

build-image:
	@docker build -t ${APP_NAME} .
	@docker image prune -f

lint-image:
	@docker run --rm -i hadolint/hadolint < Dockerfile

clean-image:
	@docker rmi -f ${APP_NAME}
