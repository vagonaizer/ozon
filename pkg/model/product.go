package model

type GetProductResponse struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

type GetProductRequest struct {
	Token string `json:"token"`
	SkuID uint32 `json:"sku"`
}
