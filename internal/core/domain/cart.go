package domain

type CartItem struct {
	SKU   int64
	Name  string
	Count uint16
	Price uint32
}

type Cart struct {
	UserID     int64
	Items      []CartItem
	TotalPrice uint32
}
