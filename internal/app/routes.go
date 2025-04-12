package app

import "net/http"

func (app *App) RegisterRoutes() {
	http.HandleFunc("POST /user/{user_id}/cart/{sku_id}", app.Handler.AddProduct)
	http.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", app.Handler.RemoveProduct)
	http.HandleFunc("DELETE /user/{user_id}/cart", app.Handler.ClearCart)
	http.HandleFunc("GET /user/{user_id}/cart", app.Handler.GetCart)
}
