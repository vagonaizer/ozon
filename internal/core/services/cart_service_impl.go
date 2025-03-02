package services

import (
	"errors"

	"github.com/vagonaizer/go-cart/internal/core/domain"
	"github.com/vagonaizer/go-cart/internal/core/ports"
)

type CartServiceImpl struct {
	productService ports.ProductService
	carts          map[int64]*domain.Cart
}

func NewCartService(productService ports.ProductService) *CartServiceImpl {
	return &CartServiceImpl{
		productService: productService,
		carts:          make(map[int64]*domain.Cart),
	}
}

func (s *CartServiceImpl) AddItem(userID int64, skuID int64, count uint16) (*domain.Cart, error) {
	product, err := s.productService.GetProduct(skuID)
	if err != nil {
		return nil, err
	}

	cart, exists := s.carts[userID]
	if !exists {
		cart = &domain.Cart{
			UserID: userID,
			Items:  []domain.CartItem{},
		}
		s.carts[userID] = cart
	}

	for i, item := range cart.Items {
		if item.SKU == skuID {
			cart.Items[i].Count += count
			cart.TotalPrice += product.Price * uint32(count)
			return cart, nil
		}
	}

	cart.Items = append(cart.Items, domain.CartItem{
		SKU:   skuID,
		Name:  product.Name,
		Count: count,
		Price: product.Price,
	})
	cart.TotalPrice += product.Price * uint32(count)

	return cart, nil
}

func (s *CartServiceImpl) RemoveItem(userID int64, skuID int64) error {
	cart, exists := s.carts[userID]
	if !exists {
		return nil
	}

	for i, item := range cart.Items {
		if item.SKU == skuID {
			cart.TotalPrice -= item.Price * uint32(item.Count)
			cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			break
		}
	}

	return nil
}

func (s *CartServiceImpl) ClearCart(userID int64) error {
	delete(s.carts, userID)
	return nil
}

func (s *CartServiceImpl) GetCart(userID int64) (*domain.Cart, error) {
	cart, exists := s.carts[userID]
	if !exists || len(cart.Items) == 0 {
		return nil, errors.New("cart not found or empty")
	}

	return cart, nil
}
