package storage

type Cart struct {
	UserID     int64               `json:"user_id"`
	TotalPrice uint32              `json:"total_price"`
	Items      map[int64]CartItems `json:"items"`
}

type CartItems struct {
	SKUID int64  `json:"sku_id"`
	Name  string `json:"item_name"`
	Price uint32 `json:"item_prince"`
	Count uint32 `json:"item_count"`
}
