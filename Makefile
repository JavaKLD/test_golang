include .env
export

GENERATOR=openapitools/openapi-generator-cli
GEN_LANG=go-server
GEN_OUTPUT=gen
OPENAPI_FILE=openapi.yml
PROJECT_DIR := $(shell pwd)
GO:=go
COVERAGE_FILE := coverage.out

TEST=go test
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html
TEST_DOCKER_COMPOSE := docker compose --file docker-compose.yml

DB_URL=mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)


run:
	go run cmd/appname/main.go

lint:
	golangci-lint run --config .golangci.yml

#Запуск всех тестов
test:
	go test -coverpkg=dolittle2/internal/utils -coverprofile=$(COVERAGE_FILE) ./tests
	go test -coverpkg=dolittle2/internal/server -coverprofile=$(COVERAGE_FILE) ./tests

#Только unit
test-unit:
	go test -coverpkg=dolittle2/internal/utils -coverprofile=$(COVERAGE_FILE) ./tests

#Только интеграционные тесты
test-integration:
	go test -coverpkg=dolittle2/internal/server -coverprofile=$(COVERAGE_FILE) ./tests

#Миграции
migrate-up:
	migrate -path ./database/migrations -database "$(DB_URL)" up

#Откат миграции
migrate-down:
	migrate -path ./database/migrations -database "$(DB_URL)" down

test-infrastructure:
	$(TEST_DOCKER_COMPOSE) up --detach --build
	$(TEST_DOCKER_COMPOSE) logs --follow

test-infrastructure-down:
	$(TEST_DOCKER_COMPOSE) down --remove-orphans
#
docker-migrate-up:
	docker run --rm \
      -e GOOSE_DRIVER=mysql \
      -e GOOSE_DBSTRING="root:strong_password@tcp(mysql:3306)/app_db?parseTime=true" \
      --network=dolittle2_default \
      goose-migration

#Загрузка зависимостей
install-deps:
	go install go install github.com/dkorunic/betteralign/cmd/betteralign@latest
	go install github.com/google/uuid v1.6.0
	go install github.com/joho/godotenv v1.5.1
	go install github.com/labstack/echo/v4 v4.13.3
	go install github.com/stretchr/testify v1.10.0
	go install google.golang.org/grpc v1.72.0
	go install google.golang.org/protobuf v1.36.6
	go install gorm.io/driver/mysql v1.5.7
	go install gorm.io/gorm v1.26.0

#Линтер
betteralign:
	$(GOPATH)/bin/betteralign ./...

#Docker

#OpenAPI
generate:
	docker run --rm -v $(PROJECT_DIR):/local $(GENERATOR) generate \
	  -i /local/$(OPENAPI_FILE) \
	  -g $(GEN_LANG) \
	  -o /local/$(GEN_OUTPUT)

#ProtoGEN
protoc:
	protoc --go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		proto/schedule.proto

clean-openapi:
	rm -rf $(GEN_OUTPUT)

clean-proto:
	rm -f proto/*.pb.go