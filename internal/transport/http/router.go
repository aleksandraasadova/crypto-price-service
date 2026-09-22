package transport

import (
	"net/http"

	"github.com/aleksandraasadova/crypto-price-service/internal/transport/middleware"
)

type RouterDeps struct {
	AuthService   AuthService
	CryptoService CryptoService
}

func NewRouter(r *RouterDeps) *http.ServeMux {
	m := http.NewServeMux()

	authHandler := NewAuthHandler(r.AuthService)
	cryptoHandler := NewCryptoHandler(r.CryptoService)

	m.HandleFunc("POST /auth/register", authHandler.Register)
	m.HandleFunc("POST /auth/login", authHandler.Login)
	m.HandleFunc("GET /crypto", middleware.JWT(cryptoHandler.GetCryptos))
	m.HandleFunc("POST /crypto", middleware.JWT(cryptoHandler.AddCrypto))
	m.HandleFunc("GET /crypto/{symbol}", middleware.JWT(cryptoHandler.GetCrypto))
	m.HandleFunc("PUT /crypto/{symbol}/refresh", middleware.JWT(cryptoHandler.RefreshCrypto))
	m.HandleFunc("GET /crypto/{symbol}/history", middleware.JWT(cryptoHandler.GetCryptoHistory))
	m.HandleFunc("GET /crypto/{symbol}/stats", middleware.JWT(cryptoHandler.GetCryptoStats))
	m.HandleFunc("DELETE /crypto/{symbol}", middleware.JWT(cryptoHandler.DeleteCrypto))

	return m
}
