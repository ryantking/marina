PWD := $(shell pwd)
GOPATH := $(shell go env GOPATH)

all: deps build run

deps:
	@GO111MODULE=on go get -u github.com/myitcv/gobin
	@gobin github.com/codegangsta/gin@cafe2ce98974a3dcca6b92ce393a91a0b58b8133
	@gobin github.com/maxcnunes/waitforit@v2.4.1
	@gobin github.com/golangci/golangci-lint/cmd/golangci-lint@v1.16.0

.PHONY: vendor
vendor:
	@GO111MODULE=on go mod tidy
	@GO111MODULE=on go mod download

build:
	@echo "Building Marina Server to $(PWD)/marinad ..."
	@GO111MODULE=on go build -o $(PWD)/marinad github.com/ryantking/marina/cmd/marinad
run:
	@$(PWD)/marinad

test: lint
	@GO111MODULE=on go test -race -covermode=atomic -coverprofile=coverage.txt github.com/ryantking/marina/pkg/...

lint:
	@echo "Running golangci-lint"
	@golangci-lint run ./pkg/...

up:
	@docker-compose up -d

down:
	@docker-compose down

gin:
	@echo "Starting local development server with gin"
	@waitforit -file=$(PWD)/waitforit.json
	@GIN_PORT=8081 BIN_APP_PORT=8080 GIN_PATH=$(PWD) GIN_BUILD=$(PWD)/cmd/marinad gin run cmd/marinad/main.go

clean:
	@echo "Cleaning up all generated files"
	@find . -name '*.test' | xargs rm -fv
	@rm -fv $(PWD)/gin-bin
	@rm -fv $(PWD)/coverage.txt
	@rm -fv $(PWD)/marinad

