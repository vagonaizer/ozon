package repository

import "testing"

func BenchmarkAddItem(b *testing.B) {
	repo := NewCartRepository()

	for i := 0; i < b.N; i++ {
		_ = repo.AddItem(1, int64(i), 1)
	}
}

func TestCartRepo_AddItem(t *testing.T) {
	repo := NewCartRepository()
	err := repo.AddItem(1, 1001, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cart, err := repo.GetCart(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cart.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(cart.Items))
	}

	if cart.Items[0].SkuID != 1001 || cart.Items[0].Count != 2 {
		t.Errorf("unexpected item data: %+v", cart.Items[0])
	}
}

func TestCartRepo_AddItem_Increment(t *testing.T) {
	repo := NewCartRepository()
	repo.AddItem(1, 1001, 2)
	repo.AddItem(1, 1001, 3)

	cart, _ := repo.GetCart(1)

	if cart.Items[0].Count != 5 {
		t.Errorf("expected count 5, got %d", cart.Items[0].Count)
	}
}

func TestCartRepo_RemoveItem(t *testing.T) {
	repo := NewCartRepository()
	repo.AddItem(1, 1001, 1)
	repo.RemoveItem(1, 1001)

	cart, err := repo.GetCart(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cart.Items) != 0 {
		t.Errorf("expected 0 items, got %d", len(cart.Items))
	}
}

func TestCartRepo_ClearCart(t *testing.T) {
	repo := NewCartRepository()
	repo.AddItem(1, 1001, 1)
	repo.AddItem(1, 1002, 2)

	err := repo.ClearCart(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = repo.GetCart(1)
	if err == nil {
		t.Fatal("expected error when getting cleared cart")
	}
}

func TestCartRepo_GetCart_NotFound(t *testing.T) {
	repo := NewCartRepository()

	_, err := repo.GetCart(99)
	if err == nil {
		t.Fatal("expected ErrCartNotFound, got nil")
	}
}
