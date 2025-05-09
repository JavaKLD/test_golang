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

#Загрузка зависимостей
install-deps:
	go install github.com/google/uuid v1.6.0
	go install github.com/joho/godotenv v1.5.1
	go install github.com/labstack/echo/v4 v4.13.3
	go install github.com/stretchr/testify v1.10.0
	go install google.golang.org/grpc v1.72.0
	go install google.golang.org/protobuf v1.36.6
	go install gorm.io/driver/mysql v1.5.7
	go install gorm.io/gorm v1.26.0

betteralign:
	betteralign -apply ./...


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