package product

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"route256/cart/pkg/middleware/logging"
	"route256/cart/pkg/model"
)

type (
	GPResponse = model.GetProductResponse
	GPRequest  = model.GetProductRequest
)

type ProductService interface {
	GetProduct(skuID uint32) (*GPResponse, error)
}

type ProductClient struct {
	client  *http.Client
	baseURL string // должен считываться с config.yml (потом прикрутить или хз НЕ ЗАБЫТЬ)
	token   string
}

func NewProductClient(client *http.Client, baseURL string, token string) *ProductClient {
	return &ProductClient{
		client:  client,
		baseURL: baseURL,
		token:   token,
	}
}

func (pc *ProductClient) GetProduct(skuID uint32) (*GPResponse, error) {
	logger := logging.GetLogger()
	logger.Info("ProductClient.GetProduct called with skuID=%d", skuID)

	reqBody := GPRequest{
		Token: pc.token,
		SkuID: skuID,
	}
	logger.Debug("Constructed GPRequest: %+v", reqBody)

	body, err := json.Marshal(reqBody)
	if err != nil {
		logger.Error("Error marshaling request: %v", err)
		return nil, fmt.Errorf("marshal request: %w", err)
	}
	logger.Debug("Marshalled request body: %s", string(body))

	url := pc.baseURL + "/get_product"
	logger.Info("POST request to URL: %s", url)
	resp, err := pc.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		logger.Error("HTTP post error: %v", err)
		return nil, fmt.Errorf("http post: %w", err)
	}
	defer resp.Body.Close()
	logger.Info("Received HTTP response with status: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		logger.Error("Unexpected status code: %d", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result GPResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		logger.Error("Error decoding response: %v", err)
		return nil, fmt.Errorf("decode response: %w", err)
	}
	logger.Info("Decoded response: %+v", result)

	return &result, nil
}
