package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/vagonaizer/go-cart/internal/core/ports"
)

const (
	StatusEnhanceYourCalm = 420
)

type ProductRepository struct {
	baseURL string
	client  *http.Client
}

func NewProductRepository(baseURL string) *ProductRepository {
	return &ProductRepository{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (r *ProductRepository) GetProduct(skuID int64) (*ports.Product, error) {
	requestBody := map[string]interface{}{
		"token": "testtoken",
		"sku":   skuID,
	}
	requestBodyBytes, _ := json.Marshal(requestBody)

	var productResponse struct {
		Name  string `json:"name"`
		Price uint32 `json:"price"`
	}

	// Ретраи
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		resp, err := r.client.Post(r.baseURL+"/get_product", "application/json", bytes.NewReader(requestBodyBytes))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusOK:
			if err := json.NewDecoder(resp.Body).Decode(&productResponse); err != nil {
				return nil, err
			}
			return &ports.Product{
				Name:  productResponse.Name,
				Price: productResponse.Price,
			}, nil
		case StatusEnhanceYourCalm, http.StatusTooManyRequests:
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		default:
			return nil, errors.New("product not found")
		}
	}

	return nil, errors.New("max retries exceeded")
}
