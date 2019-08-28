SHELL=/bin/bash
ROOT_DIR := $(shell pwd)
IMAGE_TAG := $(shell git rev-parse HEAD)
IMAGE_NAME := "quay.io/companyname/dummy_project"
REGISTRY := blabla.dkr.ecr.us-west-2.amazonaws.com

.PHONY: ci
ci: deps deps_check lint build test

.PHONY: deps
deps:
	dep ensure -v

.PHONY: deps_check
deps_check:
	@test -z "$(shell git status -s ./vendor ./Gopkg.*)"

.PHONY: grpcgen
grpcgen:
	go build -o protoc-gen-go ./vendor/github.com/golang/protobuf/protoc-gen-go
	protoc -I ./vendor -I api api/service.proto --plugin=./protoc-gen-go --go_out=plugins=grpc:api
	go build ./api

.PHONY: build
build:
	go build -o artifacts/svc

.PHONY: run
run:
	go run ./main.go

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -cover -v `go list ./...`

.PHONY: dockerise
dockerise:
	docker build -t "${IMAGE_NAME}:${IMAGE_TAG}" -f Dockerfile .

.PHONY: deploy
deploy:
	`AWS_SHARED_CREDENTIALS_FILE=~/.aws/credentials AWS_PROFILE=iam aws ecr get-login --region us-west-2 --no-include-email`
	docker push ${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}
	docker tag ${IMAGE_NAME}:${IMAGE_TAG} ${REGISTRY}/${IMAGE_NAME}:latest
	docker push ${REGISTRY}/${IMAGE_NAME}:latest


.PHONY: mock_service_db
mock_service_db:
	mockgen -source=service/storage.go -destination=service/mock/storage.go

.PHONY: run_postgresql
run_postgresql:
	docker run -d --name dummy_postgresql -v ${ROOT_DIR}/tmp/sql/data:/var/lib/postgresql/data -p 5432:5432 postgres:11

.PHONY: run_redis
run_redis:
	docker run --name dummy_redis -p 6379:6379 -d redis

.PHONY: exec_redis_sh
exec_redis_sh:
	docker exec -it dummy_redis sh
    # redis-cli