package storage

// Storage хранит корзины пользователей в памяти.
type Storage struct {
	carts map[int64]*Cart
}

// NewStorage создаёт новое хранилище корзин.
func NewStorage() *Storage {
	return &Storage{
		carts: make(map[int64]*Cart),
	}
}

// AddItem добавляет товар в корзину пользователя.
// Если товар уже существует, увеличивается его количество.
func (s *Storage) AddItem(userID int64, item CartItems) {
	cart, exists := s.carts[userID]
	if !exists {
		cart = &Cart{
			UserID: userID,
			Items:  make(map[int64]CartItems),
		}
		s.carts[userID] = cart
	}

	if existingItem, found := cart.Items[item.SKUID]; found {
		existingItem.Count += item.Count
		cart.Items[item.SKUID] = existingItem
	} else {
		cart.Items[item.SKUID] = item
	}

	cart.TotalPrice += item.Price * item.Count
}

// DeleteItem удаляет товар из корзины пользователя по sku_id.
func (s *Storage) DeleteItem(userID int64, skuID int64) {
	cart, exists := s.carts[userID]
	if !exists {
		return
	}

	if item, found := cart.Items[skuID]; found {
		cart.TotalPrice -= item.Price * item.Count
		delete(cart.Items, skuID)
	}
}

// ClearItems очищает корзину пользователя.
func (s *Storage) ClearItems(userID int64) {
	cart, exists := s.carts[userID]
	if !exists {
		return
	}

	cart.Items = make(map[int64]CartItems)
	cart.TotalPrice = 0
}

// GetItems возвращает все товары из корзины пользователя.
func (s *Storage) GetItems(userID int64) []CartItems {
	cart, exists := s.carts[userID]
	if !exists {
		return []CartItems{}
	}

	items := make([]CartItems, 0, len(cart.Items))
	for _, item := range cart.Items {
		items = append(items, item)
	}

	return items
}
