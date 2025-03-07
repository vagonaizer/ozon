package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/vagonaizer/go-cart/internal/core/ports"
)

type CartHandler struct {
	cartService ports.CartService
}

func NewCartHandler(cartService ports.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

func (h *CartHandler) AddItem(w http.ResponseWriter, r *http.Request) {

	// userID валидация

	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil || userID <= 0 {
		http.Error(w, "invalid user ID: user ID must be a positive integer", http.StatusBadRequest)
		return
	}

	skuID, err := strconv.ParseInt(r.PathValue("sku_id"), 10, 64)
	if err != nil || skuID <= 0 {
		http.Error(w, "invalid SKU ID: SKU ID must be a positive integer", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		Count uint16 `json:"count"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest) // исправить ошибки
		return
	}

	if requestBody.Count <= 0 {
		http.Error(w, "invalid count: count must be a non-negative integer", http.StatusBadRequest)
		return
	}

	if err := h.cartService.AddItem(userID, skuID, requestBody.Count); err != nil {
		http.Error(w, err.Error(), http.StatusPreconditionFailed)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *CartHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	skuID, _ := strconv.ParseInt(r.PathValue("sku_id"), 10, 64)

	if err := h.cartService.RemoveItem(userID, skuID); err != nil {
		http.Error(w, err.Error(), http.StatusPreconditionFailed) // везде исправить
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.ParseInt(r.PathValue("user_id"), 10, 64)

	if err := h.cartService.ClearCart(userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil || userID <= 0 {
		http.Error(w, "invalid user ID: user ID must be a positive integer", http.StatusBadRequest)
		return
	}

	cart, err := h.cartService.GetCart(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cart); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
