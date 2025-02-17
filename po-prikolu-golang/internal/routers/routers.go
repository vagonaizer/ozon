package router

import (
	"fmt"
	"net/http"
	"project/internal/handlers" // Обновите путь импорта согласно вашему проекту

	"github.com/gorilla/mux"
)

func TestFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Println("lalla")
	return
}

// NewRouter создаёт и возвращает настроенный маршрутизатор.
func NewRouter(cartHandler *handlers.CartHandler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/user/", TestFunc).Methods("GET")
	// Добавление товара в корзину: POST /user/{user_id}/cart/{sku_id}
	r.HandleFunc("/user/{user_id}/cart/{sku_id}", cartHandler.AddItem).Methods("POST")

	// Удаление товара из корзины: DELETE /user/{user_id}/cart/{sku_id}
	r.HandleFunc("/user/{user_id}/cart/{sku_id}", cartHandler.DeleteItem).Methods("DELETE")

	// Очистка корзины: DELETE /user/{user_id}/cart
	r.HandleFunc("/user/{user_id}/cart", cartHandler.ClearItems).Methods("DELETE")

	// Получение корзины: GET /user/{user_id}/cart
	r.HandleFunc("/user/{user_id}/cart", cartHandler.GetItems).Methods("GET")

	return r
}
