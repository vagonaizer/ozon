package main

import (
	"log"
	"net/http"
	"os"

	"github.com/vagonaizer/go-cart/internal/adapters/handler"
	"github.com/vagonaizer/go-cart/internal/adapters/repository"
	"github.com/vagonaizer/go-cart/internal/core/services"
	"github.com/vagonaizer/go-cart/internal/logger"
	"gopkg.in/yaml.v3"
)

// router
// валидация
//

type Config struct {
	Port           string `yaml:"port"`
	ProductService string `yaml:"product_service"`
}

func main() {

	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Ошибка чтения файла: %v", err)
	}

	var config Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Ошибка парсинга YAML: %v", err)
	}

	productRepo := repository.NewProductRepository(config.ProductService)
	cartService := services.NewCartService(productRepo)
	cartHandler := handler.NewCartHandler(cartService)

	logger := logger.New()

	http.Handle("POST /user/{user_id}/cart/{sku_id}", handler.LoggingMiddleware(http.HandlerFunc(cartHandler.AddItem), logger))
	http.Handle("DELETE /user/{user_id}/cart/{sku_id}", handler.LoggingMiddleware(http.HandlerFunc(cartHandler.RemoveItem), logger))
	http.Handle("DELETE /user/{user_id}/cart", handler.LoggingMiddleware(http.HandlerFunc(cartHandler.ClearCart), logger))
	http.Handle("GET /user/{user_id}/cart", handler.LoggingMiddleware(http.HandlerFunc(cartHandler.GetCart), logger))

	log.Println("Server started at :8082")
	if err := http.ListenAndServe(config.Port, nil); err != nil {
		log.Fatal(err)
	}
}
