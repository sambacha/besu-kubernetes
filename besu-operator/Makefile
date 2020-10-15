
# Makefile for besu Operator

.PHONY: all builder builder-image image package local run fmt lint test cluster-up cluster-down int-test

all: image

builder:
	@echo Creating container image to build besu-operator
	@docker build -f ./build/Dockerfile.builder -t besu/besu-operator-builder .

builder-image:
	@echo Using builder container to build besu-operator
	@mkdir -p ./build/_output/bin
	@docker run -v /var/run/docker.sock:/var/run/docker.sock -v ${PWD}:/opt/app-root/src/besu-operator -w /opt/app-root/src/besu-operator -u root -it besu/besu-operator-builder bash -c "operator-sdk build --verbose besu/besu-operator"

image:
	@echo Building besu-operator image
	@operator-sdk build --verbose besu/besu-operator

local:
	@echo Building besu-operator-local binary only
	@mkdir -p ./build/_output/bin
	@go build -v -o ./build/_output/bin/besu-operator-local ./cmd/manager

generate:
	@echo Running operator-sdk generate k8s
	@operator-sdk generate k8s
	@echo Running operator-sdk generate crds
	@operator-sdk generate crds
	@rm -f deploy/crds/*_cr.yaml
	@build/make_bundle.sh

package: lint fmt generate image
	@build/package.sh

run:
	@OPERATOR_NAME=besu-operator operator-sdk run --local

fmt:
	@gofmt -l -w `find ./ -name "*.go"`

lint:
	@golint ./...

cluster-up:
	@test/deploy-cluster.sh up

cluster-down:
	@test/deploy-cluster.sh down

int-test:
	@echo Run integration test
	@test/run-tests.sh
