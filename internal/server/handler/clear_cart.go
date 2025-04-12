package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"route256/cart/pkg/utils"
)

func (h *Handler) ClearCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userIDStr := r.PathValue("user_id")
	if userIDStr == "" {
		http.Error(w, "missing user_id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	if err := h.cartService.ClearCart(userID); err != nil {
		http.Error(w, fmt.Sprintf("cartService.ClearCart: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	utils.SuccessReponse(w)
}
