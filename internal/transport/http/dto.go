package transport

import (
	"time"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type RegisterUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterUserResponse struct {
	Token string `json:"token"`
}

type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	Token string `json:"token"`
}

type CryptoDTO struct {
	Symbol       string    `json:"symbol"`
	Name         string    `json:"name"`
	CurrentPrice float64   `json:"current_price"`
	LastUpdated  time.Time `json:"last_updated"`
}

type GetCryptosResponse struct {
	Cryptos []*CryptoDTO `json:"cryptos"`
}

type AddCryptoRequest struct {
	Symbol string `json:"symbol"`
}

type AddCryptoResponse struct {
	Crypto *CryptoDTO `json:"crypto"`
}

type GetCryptoResponse struct {
	Crypto *CryptoDTO
}

type RefreshCryptoResponse struct {
	Crypto *CryptoDTO `json:"crypto"`
}

type CryptoHistoryDTO struct {
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}

type GetHistory struct {
	Symbol  string              `json:"symbol"`
	History []*CryptoHistoryDTO `json:"history"`
}

type StatsDTO struct {
	MinPrice       float64 `json:"min_price"`
	MaxPrice       float64 `json:"max_price"`
	AvgPrice       float64 `json:"avg_price"`
	PriceChange    float64 `json:"price_change"`
	PriceChangePct float64 `json:"price_change_percent"`
	RecordsCount   int     `json:"records_count"`
}

type GetCryptoStatsResponse struct {
	Symbol       string    `json:"symbol"`
	CurrentPrice float64   `json:"current_price"`
	Stats        *StatsDTO `json:"stats"`
}
