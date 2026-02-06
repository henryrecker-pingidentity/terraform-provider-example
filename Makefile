SHELL := /bin/bash

.PHONY: install generate fmt vet test testacc golangcilint tfproviderlint tflint terrafmtlint importfmtlint devcheck devchecknotest

default: install

install:
	go mod tidy
	go install .

generate:
	go generate ./...
	go fmt ./...
	go vet ./...

fmt:
	go fmt ./...

vet:
	go vet ./...
	
test:
	go test -parallel=4 ./...

testacc:
	TF_ACC=1 go test ./... -timeout 10m -v -p 4 --count=1

devchecknotest: install golangcilint generate

devcheck: devchecknotest testacc

golangcilint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run --timeout 5m ./internal/...
