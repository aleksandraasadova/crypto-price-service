package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"

	"github.com/aleksandraasadova/crypto-price-service/internal/clients/coingecko"
	"github.com/aleksandraasadova/crypto-price-service/internal/repository/memory"
	"github.com/aleksandraasadova/crypto-price-service/internal/service"
	transport "github.com/aleksandraasadova/crypto-price-service/internal/transport/http"
	"github.com/aleksandraasadova/crypto-price-service/internal/transport/server"
)

const addr = ":8080"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("FATAL: failed to load .env file: %v", err)
	}

	// Repositories
	userRepo := memory.NewUserRepo()
	cryptoRepo := memory.NewCryptoRepo()

	//Clients
	cgClient := coingecko.NewClient()

	// Services
	authService := service.NewAuthService(userRepo, os.Getenv("JWT"))
	cryptoService := service.NewCryptoService(cryptoRepo, cgClient)

	mux := transport.NewRouter(&transport.RouterDeps{
		AuthService:   authService,
		CryptoService: cryptoService,
	})

	srv := server.New(addr, mux)
	if err := srv.Start(); err != nil {
		slog.Error("server error", "err", err)
	}
}
