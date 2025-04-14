package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"route256/cart/internal/service/cart"
	"route256/cart/pkg/model"
	"route256/cart/pkg/utils"

	"github.com/go-playground/validator/v10"
)

type cartService interface {
	AddItem(userID int64, skuID int64, count uint16) error
	RemoveItem(userID int64, skuID int64) error
	ClearCart(userID int64) error
	GetCart(userID int64) (*model.Cart, error)
}

type Handler struct {
	cartService cartService
}

func NewHandler(cs cartService) *Handler {
	return &Handler{cartService: cs}
}

type AddProductRequest struct {
	Count uint16 `json:"count" validate:"gt=0"`
}

func (h *Handler) AddProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userIDStr := r.PathValue("user_id")
	skuIDStr := r.PathValue("sku_id")
	if userIDStr == "" || skuIDStr == "" {
		http.Error(w, "{}", http.StatusBadRequest)
		return
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "{}", http.StatusBadRequest)
		return
	}
	skuID, err := strconv.ParseInt(skuIDStr, 10, 64)
	if err != nil {
		http.Error(w, "{}", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "{}", http.StatusInternalServerError)
		return
	}
	var req AddProductRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "{}", http.StatusBadRequest)
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(req); err != nil {
		http.Error(w, "{}", http.StatusBadRequest)
		return
	}

	if err := h.cartService.AddItem(userID, skuID, req.Count); err != nil {
		// Маппинг ошибок на статус-коды:
		switch err {
		case cart.ErrInvalidUserID, cart.ErrInvalidSkuID, cart.ErrInvalidProductCount:
			http.Error(w, "{}", http.StatusBadRequest)
		case cart.ErrProductServiceIssue:
			http.Error(w, "{}", http.StatusPreconditionFailed)
		default:
			http.Error(w, "{}", http.StatusInternalServerError)
		}
		return
	}

	utils.SuccessResponse(w)
}
