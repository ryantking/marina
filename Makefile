PWD := $(shell pwd)
GOPATH := $(shell go env GOPATH)

all: deps build run

deps:
	@GO111MODULE=on go get -u github.com/myitcv/gobin
	@GO111MODULE=off go get -tags 'mysql' -u github.com/golang-migrate/migrate/cmd/migrate
	@gobin github.com/codegangsta/gin@cafe2ce98974a3dcca6b92ce393a91a0b58b8133
	@gobin github.com/jwilder/dockerize@v0.6.1
	@gobin github.com/golangci/golangci-lint/cmd/golangci-lint@v1.16.0

.PHONY: vendor
vendor:
	@GO111MODULE=on go mod tidy
	@GO111MODULE=on go mod vendor

build:
	@echo "Building Marina Server to $(PWD)/marinad ..."
	@GO111MODULE=on go build -mod=vendor -o $(PWD)/marinad github.com/ryantking/marina/cmd/marinad

run:
	@$(PWD)/marinad

test: lint
	@GO111MODULE=on go test -mod=vendor -race -covermode=atomic -coverprofile=coverage.txt github.com/ryantking/marina/pkg/...

lint:
	@echo "Running golangci-lint"
	@GO111MODULE=on golangci-lint run ./...

up:
	@docker-compose up -d
	@dockerize -wait tcp://127.0.0.1:3306 -timeout 120s

down:
	@docker-compose down

migrate_up:
	@migrate -path "./migrations" -database "mysql://marina:marina@tcp(localhost:3306)/marinatest" up

migrate_down:
	@migrate -path "./migrations" -database "mysql://marina:marina@tcp(localhost:3306)/marinatest" down

destroydb:
	@migrate -path "./migrations" -database "mysql://marina:marina@tcp(localhost:3306)/marinatest" force 1
	@migrate -path "./migrations" -database "mysql://marina:marina@tcp(localhost:3306)/marinatest" down

bucket:
	@docker-compose run minio_bucket

gin:
	@echo "Starting local development server with gin"
	@GIN_PORT=8081 BIN_APP_PORT=8080 GIN_PATH=$(PWD) GIN_BUILD=$(PWD)/cmd/marinad gin run cmd/marinad/main.go

clean:
	@echo "Cleaning up all generated files"
	@find . -name '*.test' | xargs rm -fv
	@rm -fv $(PWD)/gin-bin
	@rm -fv $(PWD)/coverage.txt
	@rm -fv $(PWD)/marinad

