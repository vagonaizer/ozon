// - хз может валидацию добавить или я хз)
package repository

import (
	"log"
	"route256/cart/pkg/model"
)

type CartRepository interface {
	AddItem(userID int64, skuID int64, count uint16) error
	RemoveItem(userID int64, skuID int64) error
	ClearCart(userID int64) error
	GetCart(userID int64) (*model.Cart, error)
}

type CartRepo struct {
	// UserID -> SkuID -> Item
	items map[int64]map[int64]model.Item
}

func NewCartRepository() *CartRepo {
	return &CartRepo{
		items: make(map[int64]map[int64]model.Item),
	}
}

func (r *CartRepo) AddItem(userID int64, skuID int64, count uint16) error {
	// проверка
	// существует ли корзина по userid
	if _, ok := r.items[userID]; !ok {
		log.Printf("(AddItem) cart_repo.go: New cart was created for UserID = %d", userID)
		r.items[userID] = make(map[int64]model.Item)
	}

	item, exists := r.items[userID][skuID]
	if exists {
		item.Count += count
	} else {
		item = model.Item{
			SkuID: skuID,
			Count: count,
		}
	}

	r.items[userID][skuID] = item

	return nil
}

func (r *CartRepo) RemoveItem(userID int64, skuID int64) error {
	// if _, ok := r.items[userID]; !ok {
	// 	log.Printf("(RemoveItem) cart_repo.go: Cart not found or not existing, userID=%d", userID)
	// 	return ErrCartNotFound
	// } нахуй не нужно

	// существует ли айтем

	delete(r.items[userID], skuID)

	return nil
}

func (r *CartRepo) ClearCart(userID int64) error {
	if _, ok := r.items[userID]; !ok {
		log.Printf("(ClearCart) cart_repo.go: Cart not found or not existing, userID=%d", userID)
		return ErrCartNotFound
	}

	delete(r.items, userID)

	return nil
}

func (r *CartRepo) GetCart(userID int64) (*model.Cart, error) {
	userItems, exists := r.items[userID]
	if !exists {
		log.Printf("(ClearCart) cart_repo.go: Cart not found or empty, userID=%d", userID)
		return nil, ErrCartNotFound
	}

	cart := &model.Cart{
		UserID:     userID,
		TotalPrice: 0,
		Items:      make(model.Items, 0),
	}

	totalPrice := uint32(0)
	for _, item := range userItems {
		cart.Items = append(cart.Items, item)
		totalPrice += item.Price * uint32(item.Count)
	}
	cart.TotalPrice = totalPrice

	return cart, nil
}
