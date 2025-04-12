package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	// "route256/cart/pkg/model"
	// "route256/cart/pkg/utils"
)

type GetCartResponseProduct struct {
	SkuId int64  `json:"sku_id"`
	Name  string `json:"name"`
	Count uint16 `json:"count"`
	Price uint32 `json:"price"`
}

type GetCartResponse struct {
	Items      []GetCartResponseProduct `json:"items"`
	TotalPrice uint32                   `json:"total_price"`
}

func (h *Handler) GetCart(w http.ResponseWriter, r *http.Request) {
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

	cart, err := h.cartService.GetCart(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("cartService.GetCart: %v", err), http.StatusInternalServerError)
		return
	}

	if cart == nil || len(cart.Items) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{}"))
		return
	}

	items := make([]GetCartResponseProduct, 0, len(cart.Items))
	for _, item := range cart.Items {
		items = append(items, GetCartResponseProduct{
			SkuId: item.SkuID,
			Name:  item.Name,
			Count: item.Count,
			Price: item.Price,
		})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].SkuId < items[j].SkuId
	})

	response := GetCartResponse{
		Items:      items,
		TotalPrice: cart.TotalPrice,
	}

	data, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("json.Marshal: %v", err), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}
