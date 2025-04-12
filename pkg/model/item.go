// domain (entity) - Item
package model

type Item struct {
	SkuID int64  `json:"sku_id"`
	Name  string `json:"name"`
	Price uint32 `json:"price"`
	Count uint16 `json:"count"`
}

type Items []Item
