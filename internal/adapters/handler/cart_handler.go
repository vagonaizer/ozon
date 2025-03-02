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
	userID, _ := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	skuID, _ := strconv.ParseInt(r.PathValue("sku_id"), 10, 64)

	var requestBody struct {
		Count uint16 `json:"count"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	cart, err := h.cartService.AddItem(userID, skuID, requestBody.Count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

func (h *CartHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	skuID, _ := strconv.ParseInt(r.PathValue("sku_id"), 10, 64)

	if err := h.cartService.RemoveItem(userID, skuID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.ParseInt(r.PathValue("user_id"), 10, 64)

	if err := h.cartService.ClearCart(userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.ParseInt(r.PathValue("user_id"), 10, 64)

	cart, err := h.cartService.GetCart(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}
