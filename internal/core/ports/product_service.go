package ports

type ProductService interface {
	GetProduct(skuID int64) (*Product, error)
}

type Product struct {
	Name  string
	Price uint32
}
