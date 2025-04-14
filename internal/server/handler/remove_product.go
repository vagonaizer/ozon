package handler

import (
	"net/http"
	"strconv"
)

func (h *Handler) RemoveProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userIDStr := r.PathValue("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	skuIDStr := r.PathValue("sku_id")
	skuID, err := strconv.ParseInt(skuIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid sku_id", http.StatusBadRequest)
		return
	}

	if err := h.cartService.RemoveItem(userID, skuID); err != nil {
		// Упрощаем проверку ошибок, так как оригинальные ошибки не экспортируются
		if err.Error() == "invalid user_id" || err.Error() == "invalid sku_id" { // импортировать ошибки из бизнес логики
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}
