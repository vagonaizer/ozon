// domain (entity) - Cart
package model

type Cart struct {
	UserID     int64  `json:"user_id"`
	TotalPrice uint32 `json:"total_price"`
	Items      Items  `json:"items"`
}
