package ports

type ProductService interface {
	GetProduct(skuID int64) (*Product, error)
	ListSKUs(startAfterSku int64, count int) ([]int64, error)
}

type Product struct {
	Name  string
	Price uint32
}
