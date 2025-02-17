package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"project/internal/storage" // Обновите путь импорта согласно вашей структуре проекта

	"github.com/gorilla/mux"
)

// Допустимые SKU
var validSkus = map[int64]bool{
	1076963: true,
	1148162: true,
}

// CartHandler определяет HTTP-хендлеры для работы с корзиной.
type CartHandler struct {
	Store *storage.Storage
}

// NewCartHandler создаёт новый CartHandler.
func NewCartHandler(store *storage.Storage) *CartHandler {
	return &CartHandler{Store: store}
}

// AddItem обрабатывает POST-запрос на добавление товара в корзину.
// URL: /user/{user_id}/cart/{sku_id}
// Тело запроса ожидает JSON вида: { "count": <число> }
func (h *CartHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userID, err := strconv.ParseInt(vars["user_id"], 10, 64)
	if err != nil || userID <= 0 {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	skuID, err := strconv.ParseInt(vars["sku_id"], 10, 64)
	if err != nil || skuID <= 0 {
		http.Error(w, "invalid sku", http.StatusBadRequest)
		return
	}

	// Если sku неизвестен, возвращаем 412 Precondition Failed
	if !validSkus[skuID] {
		http.Error(w, "invalid sku", http.StatusPreconditionFailed)
		return
	}

	// Декодируем тело запроса
	var payload struct {
		Count int `json:"count"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if payload.Count <= 0 {
		http.Error(w, "invalid count", http.StatusBadRequest)
		return
	}

	// Формируем объект товара.
	item := storage.CartItems{
		SKUID: skuID,
		Count: uint32(payload.Count),
	}
	// Для примера зададим цену и наименование в зависимости от sku.
	switch skuID {
	case 1076963:
		item.Price = 100 // Например, цена 100
		item.Name = "Item 1076963"
	case 1148162:
		item.Price = 200 // Например, цена 200
		item.Name = "Item 1148162"
	}

	h.Store.AddItem(userID, item)

	// Возвращаем пустой JSON с 200 OK
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

// DeleteItem обрабатывает DELETE-запрос на удаление конкретного товара.
// URL: /user/{user_id}/cart/{sku_id}
func (h *CartHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userID, err := strconv.ParseInt(vars["user_id"], 10, 64)
	if err != nil || userID <= 0 {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	skuID, err := strconv.ParseInt(vars["sku_id"], 10, 64)
	if err != nil || skuID <= 0 {
		http.Error(w, "invalid sku", http.StatusBadRequest)
		return
	}

	// Если sku неизвестен, возвращаем 412 Precondition Failed
	if !validSkus[skuID] {
		http.Error(w, "invalid sku", http.StatusPreconditionFailed)
		return
	}

	h.Store.DeleteItem(userID, skuID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

// ClearItems обрабатывает DELETE-запрос на удаление всей корзины.
// URL: /user/{user_id}/cart
func (h *CartHandler) ClearItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userID, err := strconv.ParseInt(vars["user_id"], 10, 64)
	if err != nil || userID <= 0 {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	h.Store.ClearItems(userID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

// GetItems обрабатывает GET-запрос на получение корзины пользователя.
// URL: /user/{user_id}/cart
func (h *CartHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userID, err := strconv.ParseInt(vars["user_id"], 10, 64)
	if err != nil || userID <= 0 {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	items := h.Store.GetItems(userID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}
