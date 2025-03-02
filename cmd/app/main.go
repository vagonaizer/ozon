package main

import (
	"log"
	"net/http"

	"github.com/vagonaizer/go-cart/internal/adapters/handler"
	"github.com/vagonaizer/go-cart/internal/adapters/repository"
	"github.com/vagonaizer/go-cart/internal/core/services"
)

func main() {
	productRepo := repository.NewProductRepository("http://route256.pavl.uk:8080")
	cartService := services.NewCartService(productRepo)
	cartHandler := handler.NewCartHandler(cartService)

	http.HandleFunc("POST /user/{user_id}/cart/{sku_id}", cartHandler.AddItem)
	http.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", cartHandler.RemoveItem)
	http.HandleFunc("DELETE /user/{user_id}/cart", cartHandler.ClearCart)
	http.HandleFunc("GET /user/{user_id}/cart", cartHandler.GetCart)

	log.Println("Server started at :8082")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatal(err)
	}
}
