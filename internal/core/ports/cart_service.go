package ports

import "github.com/vagonaizer/go-cart/internal/core/domain"

type CartService interface {
	AddItem(userID int64, skuID int64, count uint16) (*domain.Cart, error)
	RemoveItem(userID int64, skuID int64) error
	ClearCart(userID int64) error
	GetCart(userID int64) (*domain.Cart, error)
}
