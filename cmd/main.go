package main

// структура

// storage/storage.go - описание in_memory бд и ее методов: additem, deleteitem, clearitems, getitems
// handlers/handlers.go - хендлеры addtocart, deleteitemfromcart, clearcart, getcart
// server/server.go - маршрутизатор, обработка POST, GET, DELETE эндоопоинтов
// cmd/main.go - точка входа
// logger/logger.go - логгирование

import (
	"log"
	"net/http"

	"github.com/vagonaizer/ozon-project/internal/handlers" // Обновите путь импорта согласно вашему проекту
	router "github.com/vagonaizer/ozon-project/internal/routers"
	"github.com/vagonaizer/ozon-project/internal/storage" // Обновите путь импорта согласно вашему проекту
)

func main() {
	// Инициализация in-memory хранилища
	store := storage.NewStorage()

	// Создание экземпляра хендлера для работы с корзиной
	cartHandler := handlers.NewCartHandler(store)

	// Создание маршрутизатора
	r := router.NewRouter(cartHandler)

	log.Println(">>> Starting server on :8082 <<<") // Отладка: точно ли сервер запускается?
	err := http.ListenAndServe(":8082", r)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
