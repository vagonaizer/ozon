APP_NAME := go-route-cart
BUILD_DIR := ./bin
SOURCE_DIR := ./cmd/app

.PHONY: build run-all clean test-coverage

build:
	@echo "Сборка приложения..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(SOURCE_DIR)
	@echo "Собранный бинарник: $(BUILD_DIR)/$(APP_NAME)"

run-all: build
	@echo "Запуск приложения..."
	@$(BUILD_DIR)/$(APP_NAME)

clean:
	@echo "Очистка бинарников..."
	@rm -rf $(BUILD_DIR)
	@echo "Готово"

test-coverage:
	@echo "Запуск тестов с покрытием..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out