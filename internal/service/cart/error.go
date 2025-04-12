package cart

import "errors"

var (
	ErrInvalidUserID       = errors.New("invalid user id, should be a positive number")
	ErrInvalidSkuID        = errors.New("invalid sku id, should be a positive number")
	ErrInvalidProductCount = errors.New("invalid product count, should be a positive number")
	// ------ //
	ErrProductServiceIssue = errors.New("product service method: get_product not working")
	ErrRepositoryIssue     = errors.New("repository issue occured while trying to operate with it")
	ErrInvalidRequest      = errors.New("invalid request, check your params")
	ErrCartNotFound        = errors.New("cart not found, check your params")
)
