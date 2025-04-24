package e2e

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"route256/cart/internal/repository"
	"route256/cart/internal/server/handler"
	"route256/cart/internal/service/cart"
	"route256/cart/pkg/model"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Фейковый productService
type stubProductService struct{}

func (s *stubProductService) GetProduct(skuID uint32) (*model.GetProductResponse, error) {
	return &model.GetProductResponse{
		Name:  "TestProduct",
		Price: 150,
	}, nil
}

// Инмемори реализация репозитория
func setupServer() http.Handler {
	repo := repository.NewCartRepository()
	ps := &stubProductService{}
	svc := cart.NewCartService(repo, ps)

	// Засетапим начальные данные
	_ = svc.AddItem(1, 123, 2) // userID=1, skuID=123, count=2

	h := handler.NewHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /cart/{user_id}/list", h.GetCart)
	mux.HandleFunc("DELETE /cart/{user_id}/item/{sku_id}", h.RemoveProduct)

	return mux
}

func Test_GetCart_Positive(t *testing.T) {
	server := setupServer()

	req := httptest.NewRequest(http.MethodGet, "/cart/1/list", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)

	var parsed struct {
		Items      []handler.GetCartResponseProduct `json:"items"`
		TotalPrice uint32                           `json:"total_price"`
	}

	err := json.Unmarshal(body, &parsed)
	require.NoError(t, err)
	require.Len(t, parsed.Items, 1)

	item := parsed.Items[0]
	require.Equal(t, int64(123), item.SkuId)
	require.Equal(t, "TestProduct", item.Name)
	require.Equal(t, uint16(2), item.Count)
	require.Equal(t, uint32(150), item.Price)
	require.Equal(t, uint32(300), parsed.TotalPrice)
}

func Test_DeleteItem_Positive(t *testing.T) {
	server := setupServer()

	// DELETE /cart/1/item/123
	req := httptest.NewRequest(http.MethodDelete, "/cart/1/item/123", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// GET /cart/1/list → корзина пуста
	req = httptest.NewRequest(http.MethodGet, "/cart/1/list", nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)

	resp = w.Result()
	defer resp.Body.Close()
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	require.Equal(t, "{}", strings.TrimSpace(string(body)))
}
