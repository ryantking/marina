# Variable configuration
PWD := $(shell pwd)
GOPATH := $(shell go env GOPATH)

all: vendor build run
.PHONY: vendor up down wait-for-it migrate_up migrate_down destroydb local

# Building and running
vendor:
	@GO111MODULE=on go mod tidy
	@GO111MODULE=on go mod vendor
build:
	@echo "Building Marina Server to $(PWD)/marinad ..."
	@GO111MODULE=on go build -mod=vendor -o $(PWD)/marinad github.com/ryantking/marina/cmd/marinad
run:
	@$(PWD)/marinad

# Local dependencies
up:
	@docker-compose up -d mysql prisma minio
down:
	@docker-compose down
wait-for-it:
	@echo -n "Waiting for Minio..."
	@$(PWD)/scripts/wait-for-it.sh 'curl localhost:9000'
	@echo "done"
	@echo -n "Waiting for MySQL..."
	@$(PWD)/scripts/wait-for-it.sh 'curl localhost:3306'
	@echo "done"
	@echo -n "Waiting for Prisma..."
	@$(PWD)/scripts/wait-for-it.sh 'curl localhost:4466'
	@echo "done"

# Setup
prisma:
	@prisma deploy
bucket:
	@docker-compose run minio_bucket

# Bring up local server with hot reload
$(GOPATH)/bin/gin:
	@go get github.com/codegangsta/gin
local: $(GOPATH)/bin/gin up wait-for-it prisma bucket
	@echo "Starting local development server with gin"
	@GIN_PORT=8081 BIN_APP_PORT=8080 GIN_PATH=$(PWD) GIN_BUILD=$(PWD)/cmd/marinad gin run cmd/marinad/main.go

# Linting and Testing
$(GOPATH)/bin/golangci-lint:
	@go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
lint: $(GOPATH)/bin/golangci-lint
	@echo "Running golangci-lint"
	@GO111MODULE=on golangci-lint run ./...
test: wait-for-it prisma bucket
	@GO111MODULE=on go test -race -covermode=atomic -coverprofile=coverage.txt github.com/ryantking/marina/pkg/...

# Remove generated files
clean:
	@echo "Cleaning up all generated files"
	@find . -name '*.test' | xargs rm -fv
	@rm -fv $(PWD)/gin-bin
	@rm -fv $(PWD)/coverage.txt
	@rm -fv $(PWD)/marinad

