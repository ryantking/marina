PWD := $(shell pwd)
GOPATH := $(shell go env GOPATH)

all: deps build run

test: lint
	@GO111MODULE=on go test -race -covermode=atomic -coverprofile=coverage.txt github.com/ryantking/marina/pkg/...

.PHONY: vendor
vendor:
	@GO111MODULE=on go mod tidy
	@GO111MODULE=on go mod download

build:
	@echo "Building Marina Server to $(PWD)/marinad ..."
	@GO111MODULE=on go build -o $(PWD)/marinad github.com/ryantking/marina/cmd/marinad

run:
	@$(PWD)/marinad

deps:
	@GO111MODULE=on go get -u github.com/myitcv/gobin
	@gobin github.com/golangci/golangci-lint/cmd/golangci-lint@v1.16.0

lint:
	@echo "Running $@"
	@golangci-lint run ./pkg/...

clean:
	@echo "Cleaning up all generated files"
	@find . -name '*.test' | xargs rm -fv
	@rm $(PWD)/marinad

