SHELL        := /bin/bash -euo pipefail
APP_NAME     := aws-ecr-login-secret-updater
KUBE_LIB_VER := 1.22.1

all: build test lint

build: GOOS ?= $(shell go env GOOS)
build: GOARCH ?= $(shell go env GOARCH)
build: CGO_ENABLED ?= $(shell go env CGO_ENABLED)
build:
	GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=${CGO_ENABLED} go build -ldflags="-s -w" -trimpath -tags timetzdata -o ${APP_NAME}

test:
	@go clean -testcache
	@go test -race ./...

lint:
	@go vet ./...
	@golint -set_exit_status ./...

run: TZ ?= Asia/Tokyo
run: AWS_REGION ?= ap-northeast-1
run: AWS_ACCOUNT_ID ?= 000000000000
run: AWS_ACCESS_KEY_ID ?= AAAAAAAAAAAAAAAAAAAA
run: AWS_SECRET_ACCESS_KEY ?= AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
run: EMAIL ?= foo@example.com
run: SECRET ?= sample-ecr-login-secret
run: NAMESPACE ?= default
run:
	@TZ=${TZ} \
	AWS_REGION=${AWS_REGION} \
	AWS_ACCOUNT_ID=${AWS_ACCOUNT_ID} \
	AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
	AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
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

mod-replace-kube:
	@./go_mod_replace.sh ${KUBE_LIB_VER}
