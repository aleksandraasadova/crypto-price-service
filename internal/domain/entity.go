package domain

import (
	"errors"
	"time"
)

type User struct {
	Username string
	Password string
}

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrWrongPassword      = errors.New("wrong password")
)

type Cryptocurrency struct {
	Symbol       string
	Name         string
	CurrentPrice float64
	LastUpdated  time.Time
}

var (
	ErrCryptoNotFound = errors.New("crypto currency not found")
	ErrCryptoExists   = errors.New("crypro currency already exists")
	ErrUnknownSymbol  = errors.New("Unknown symbol")
)

type CryptoStats struct {
	Symbol         string
	CurrentPrice   float64
	MinPrice       float64
	MaxPrice       float64
	AvgPrice       float64
	PriceChange    float64
	PriceChangePct float64
	RecordsCount   int
}
