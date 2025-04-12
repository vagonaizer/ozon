//go:generate minimock -i cartRepository -o ./mocks -s "_mock.go" -g
//go:generate minimock -i productService -o ./mocks -s "_mock.go" -g
package cart

import (
	"errors"
	"log"
	"route256/cart/pkg/middleware/logging"
	"route256/cart/pkg/model"
)

type cartRepository interface {
	AddItem(userID int64, skuID int64, count uint16) error
	RemoveItem(userID int64, skuID int64) error
	ClearCart(userID int64) error
	GetCart(userID int64) (*model.Cart, error)
}

type GetProductResponse = model.GetProductResponse

type productService interface {
	GetProduct(skuID uint32) (*GetProductResponse, error)
}

type CartService struct {
	cartRepository cartRepository
	productService productService
}

func NewCartService(cartRepository cartRepository, productService productService) *CartService {
	return &CartService{
		cartRepository: cartRepository,
		productService: productService,
	}
}

func (cs *CartService) AddItem(userID int64, skuID int64, count uint16) error {
	logger := logging.GetLogger()
	logger.Info("AddItem called")

	if userID < 1 || skuID < 1 || count < 1 {
		logger.Error("Invalid params")
		return ErrInvalidRequest
	}

	if _, err := cs.productService.GetProduct(uint32(skuID)); err != nil {
		logger.Error("Product error")
		return ErrProductServiceIssue
	}

	return cs.cartRepository.AddItem(userID, skuID, count)
}

func (cs *CartService) RemoveItem(userID int64, skuID int64) error {
	logger := logging.GetLogger()
	logger.Info("RemoveItem called")

	if userID < 1 {
		logger.Error("Invalid userID")
		return ErrInvalidUserID
	}
	if skuID < 1 {
		logger.Error("Invalid skuID")
		return ErrInvalidSkuID
	}

	err := cs.cartRepository.RemoveItem(userID, skuID)
	if errors.Is(err, ErrCartNotFound) {
		return nil
	}
	if err != nil {
		logger.Error("Repository error")
		return ErrRepositoryIssue
	}

	return nil
}

func (cs *CartService) ClearCart(userID int64) error {
	if userID < 1 {
		log.Printf("(buisness-logic) cart.go - invalid userID, userID=%d", userID)
		return ErrInvalidUserID
	}
	if err := cs.cartRepository.ClearCart(userID); err != nil {
		log.Printf("(buisness-logic) cart.go - error while clearing cart, err=%v", err)
		return ErrRepositoryIssue
	}
	return nil
}

func (cs *CartService) GetCart(userID int64) (*model.Cart, error) {
	logger := logging.GetLogger()
	logger.Info("GetCart")

	cart, err := cs.cartRepository.GetCart(userID)
	if err != nil {
		logger.Error("GetCart error")
		return nil, err
	}

	for i := range cart.Items {
		product, err := cs.productService.GetProduct(uint32(cart.Items[i].SkuID))
		if err != nil {
			log.Println("GetProduct error", cart.Items[i].SkuID, err)
			return nil, err
		}
		cart.Items[i].Name = product.Name
		cart.Items[i].Price = product.Price
		cart.TotalPrice += product.Price * uint32(cart.Items[i].Count)
	}

	return cart, nil
}
