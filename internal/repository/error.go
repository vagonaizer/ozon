package repository

import "errors"

// файл с отдельными ошибками под слой репозитория

var (
	ErrCartNotFound = errors.New("cart not found or not existing")
)
