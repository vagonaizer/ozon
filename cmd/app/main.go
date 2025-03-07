package main

import (
	"log"
	"net/http"
	"os"

	"github.com/vagonaizer/go-cart/internal/adapters/handler"
	"github.com/vagonaizer/go-cart/internal/adapters/repository"
	"github.com/vagonaizer/go-cart/internal/core/services"
	"gopkg.in/yaml.v3"
)

// router
// валидация
//

type Config struct {
	Host           string `yaml:"host"`
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

	http.HandleFunc("POST /user/{user_id}/cart/{sku_id}", cartHandler.AddItem)
	http.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", cartHandler.RemoveItem)
	http.HandleFunc("DELETE /user/{user_id}/cart", cartHandler.ClearCart)
	http.HandleFunc("GET /user/{user_id}/cart", cartHandler.GetCart)

	log.Println("Server started at :8082")
	if err := http.ListenAndServe(config.Port, nil); err != nil {
		log.Fatal(err)
	}
}
