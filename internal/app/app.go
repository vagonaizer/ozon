package app

import (
	"fmt"
	"log"
	"net/http"
	"route256/cart/internal/config"
	"route256/cart/internal/repository"
	"route256/cart/internal/server/handler"
	"route256/cart/internal/service/cart"
	"route256/cart/internal/service/product"
	"route256/cart/pkg/middleware/logging"
)

type App struct {
	logger         *logging.Logger
	config         *config.Config
	productService *product.ProductClient
	cartService    *cart.CartService
	Handler        *handler.Handler
}

func NewApp() *App {

	logging.Init(logging.DEBUG)
	logger := logging.GetLogger()
	logger.Info("logger initialized")

	cfg, _ := config.GetConfig() // если не будет работать -> cleanenv
	logger.Info("config initialized")

	cartRepo := repository.NewCartRepository()
	logger.Info("repository initialized")

	productClient := product.NewProductClient(
		&http.Client{},
		cfg.PSConfig.ProductServiceURL,
		cfg.PSConfig.ProductServiceToken,
	)
	logger.Info("product client initialized")

	cartService := cart.NewCartService(cartRepo, productClient)
	logger.Info("cart service initialized")

	h := handler.NewHandler(cartService)

	return &App{
		logger:         logger,
		config:         cfg,
		productService: productClient,
		cartService:    cartService,
		Handler:        h,
	}
}

func (a *App) Run() {
	a.RegisterRoutes()

	bindIP := a.config.ServerConfig.BindIP
	if bindIP == "" {
		bindIP = "0.0.0.0"
	}
	port := a.config.ServerConfig.Port
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf("%s:%s", bindIP, port)

	a.logger.Info(fmt.Sprintf("Server listening on %s", addr))
	log.Fatal(http.ListenAndServe(addr, nil))
}
