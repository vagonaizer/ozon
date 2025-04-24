package cart_test

import (
	"errors"
	"testing"

	"route256/cart/internal/service/cart"
	"route256/cart/internal/service/cart/mocks"
	"route256/cart/pkg/model"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestCartService_AddItem(t *testing.T) {
	mc := minimock.NewController(t)
	repo := mocks.NewCartRepositoryMock(mc)
	ps := mocks.NewProductServiceMock(mc)

	svc := cart.NewCartService(repo, ps)

	t.Run("ok", func(t *testing.T) {
		ps.GetProductMock.Return(&model.GetProductResponse{Name: "Test", Price: 100}, nil)
		repo.AddItemMock.Return(nil)

		err := svc.AddItem(1, 2, 3)
		require.NoError(t, err)
	})

	t.Run("invalid input", func(t *testing.T) {
		err := svc.AddItem(0, 2, 1)
		require.ErrorIs(t, err, cart.ErrInvalidRequest)
	})

	t.Run("product service fails", func(t *testing.T) {
		ps.GetProductMock.Return(nil, errors.New("fail"))

		err := svc.AddItem(1, 2, 3)
		require.ErrorIs(t, err, cart.ErrProductServiceIssue)
	})
}

func TestCartService_RemoveItem(t *testing.T) {
	mc := minimock.NewController(t)
	repo := mocks.NewCartRepositoryMock(mc)
	ps := mocks.NewProductServiceMock(mc)

	svc := cart.NewCartService(repo, ps)

	t.Run("ok", func(t *testing.T) {
		repo.RemoveItemMock.Return(nil)

		err := svc.RemoveItem(1, 2)
		require.NoError(t, err)
	})

	t.Run("invalid userID", func(t *testing.T) {
		err := svc.RemoveItem(0, 2)
		require.ErrorIs(t, err, cart.ErrInvalidUserID)
	})

	t.Run("invalid skuID", func(t *testing.T) {
		err := svc.RemoveItem(1, 0)
		require.ErrorIs(t, err, cart.ErrInvalidSkuID)
	})

	t.Run("cart not found — пропускаем ошибку", func(t *testing.T) {
		repo.RemoveItemMock.Return(cart.ErrCartNotFound)

		err := svc.RemoveItem(1, 2)
		require.NoError(t, err)
	})

	t.Run("репозиторий ошибка", func(t *testing.T) {
		repo.RemoveItemMock.Return(errors.New("some repo err"))

		err := svc.RemoveItem(1, 2)
		require.ErrorIs(t, err, cart.ErrRepositoryIssue)
	})
}

func TestCartService_ClearCart(t *testing.T) {
	mc := minimock.NewController(t)
	repo := mocks.NewCartRepositoryMock(mc)
	ps := mocks.NewProductServiceMock(mc)

	svc := cart.NewCartService(repo, ps)

	t.Run("ok", func(t *testing.T) {
		repo.ClearCartMock.Return(nil)

		err := svc.ClearCart(1)
		require.NoError(t, err)
	})

	t.Run("invalid userID", func(t *testing.T) {
		err := svc.ClearCart(0)
		require.ErrorIs(t, err, cart.ErrInvalidUserID)
	})

	t.Run("ошибка репозитория", func(t *testing.T) {
		repo.ClearCartMock.Return(errors.New("boom"))

		err := svc.ClearCart(1)
		require.ErrorIs(t, err, cart.ErrRepositoryIssue)
	})
}

func TestCartService_GetCart(t *testing.T) {
	mc := minimock.NewController(t)
	repo := mocks.NewCartRepositoryMock(mc)
	ps := mocks.NewProductServiceMock(mc)

	svc := cart.NewCartService(repo, ps)

	t.Run("ok", func(t *testing.T) {
		repo.GetCartMock.Return(&model.Cart{
			UserID: 1,
			Items: model.Items{
				{SkuID: 101, Count: 2},
			},
		}, nil)

		ps.GetProductMock.Return(&model.GetProductResponse{Name: "Milk", Price: 100}, nil)

		cartResult, err := svc.GetCart(1)
		require.NoError(t, err)
		require.Equal(t, uint32(200), cartResult.TotalPrice)
		require.Equal(t, "Milk", cartResult.Items[0].Name)
		require.Equal(t, uint32(100), cartResult.Items[0].Price)
	})

	t.Run("ошибка GetCart", func(t *testing.T) {
		repo.GetCartMock.Return(nil, errors.New("not found"))

		c, err := svc.GetCart(1)
		require.Nil(t, c)
		require.Error(t, err)
	})

	t.Run("ошибка product service", func(t *testing.T) {
		repo.GetCartMock.Return(&model.Cart{
			UserID: 1,
			Items: model.Items{
				{SkuID: 123, Count: 1},
			},
		}, nil)

		ps.GetProductMock.Return(nil, errors.New("ps error"))

		c, err := svc.GetCart(1)
		require.Nil(t, c)
		require.Error(t, err)
	})
}
