package domain

type CartItem struct {
	SKU   int64  `json:"sku_id"`
	Name  string `json:"item_name"`
	Count uint16 `json:"item_count"`
	Price uint32 `json:"item_price"`
}

type Cart struct {
	UserID     int64      `json:"user_id"`
	Items      []CartItem `json:"items"`
	TotalPrice uint32     `json:"total_price"`
}
