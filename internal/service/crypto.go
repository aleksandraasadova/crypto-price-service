package service

import (
	"context"
	"errors"

	"github.com/aleksandraasadova/crypto-price-service/internal/clients/coingecko"
	"github.com/aleksandraasadova/crypto-price-service/internal/domain"
)

type CryptoRepository interface {
	GetAll(ctx context.Context) ([]*domain.Cryptocurrency, error)
	Get(ctx context.Context, symbol string) (*domain.Cryptocurrency, error)
	Add(ctx context.Context, crypto *domain.Cryptocurrency)
	GetHistory(ctx context.Context, symbol string) ([]*domain.Cryptocurrency, error)
	Delete(ctx context.Context, symbol string) error
}

type CryptoService struct {
	cryptoRepository CryptoRepository
	cgClient         *coingecko.Client
}

func NewCryptoService(cryptoRepository CryptoRepository, cgClient *coingecko.Client) *CryptoService {
	return &CryptoService{
		cryptoRepository: cryptoRepository,
		cgClient:         cgClient,
	}
}

func (c *CryptoService) GetCryptos(ctx context.Context) ([]*domain.Cryptocurrency, error) {
	return c.cryptoRepository.GetAll(ctx)
}

func (c *CryptoService) AddCrypto(ctx context.Context, symbol string) (*domain.Cryptocurrency, error) {
	_, err := c.cryptoRepository.Get(ctx, symbol)
	if !errors.Is(err, domain.ErrCryptoNotFound) {
		return nil, domain.ErrCryptoExists
	}

	cryptoCurrency, err := c.cgClient.GetCryptoCurrency(ctx, symbol)
	if err != nil {
		return nil, err
	}

	c.cryptoRepository.Add(ctx, cryptoCurrency)

	return cryptoCurrency, nil
}

func (c *CryptoService) GetCrypto(ctx context.Context, symbol string) (*domain.Cryptocurrency, error) {
	return c.cryptoRepository.Get(ctx, symbol)
}

func (c *CryptoService) RefreshCrypto(ctx context.Context, symbol string) (*domain.Cryptocurrency, error) {
	if _, err := c.cryptoRepository.Get(ctx, symbol); err != nil {
		return nil, err // 404
	}

	freshData, err := c.cgClient.GetCryptoCurrency(ctx, symbol)
	if err != nil {
		return nil, err
	}

	c.cryptoRepository.Add(ctx, freshData)

	return freshData, nil
}

func (c *CryptoService) GetHistory(ctx context.Context, symbol string) ([]*domain.Cryptocurrency, error) {
	return c.cryptoRepository.GetHistory(ctx, symbol)
}

func (c *CryptoService) GetStats(ctx context.Context, symbol string) (*domain.CryptoStats, error) {
	history, err := c.cryptoRepository.GetHistory(ctx, symbol)
	if err != nil {
		return nil, err
	}
	count := len(history)

	firstPrice := history[0].CurrentPrice
	lastPrice := history[count-1].CurrentPrice

	minPrice := firstPrice
	maxPrice := firstPrice
	sum := 0.0

	for _, v := range history {
		if v.CurrentPrice < minPrice {
			minPrice = v.CurrentPrice
		}
		if v.CurrentPrice > maxPrice {
			maxPrice = v.CurrentPrice
		}
		sum += v.CurrentPrice
	}

	avgPrice := sum / float64(count)
	priceChange := lastPrice - firstPrice

	var priceChangePct float64
	if firstPrice != 0 {
		priceChangePct = (priceChange / firstPrice) * 100
	}

	return &domain.CryptoStats{
		Symbol:         symbol,
		CurrentPrice:   lastPrice,
		MinPrice:       minPrice,
		MaxPrice:       maxPrice,
		AvgPrice:       avgPrice,
		PriceChange:    priceChange,
		PriceChangePct: priceChangePct,
		RecordsCount:   count,
	}, nil
}

func (c *CryptoService) DeleteCrypto(ctx context.Context, symbol string) error {
	return c.cryptoRepository.Delete(ctx, symbol)
}
