package memory

import (
	"context"
	"sync"

	"github.com/aleksandraasadova/crypto-price-service/internal/domain"
)

type CryptoRepo struct {
	mu    sync.RWMutex
	coins map[string][]*domain.Cryptocurrency // key: symbol (BTC, ETH, etc.)
}

func NewCryptoRepo() *CryptoRepo {
	return &CryptoRepo{
		coins: make(map[string][]*domain.Cryptocurrency),
	}
}

func (c *CryptoRepo) GetAll(ctx context.Context) ([]*domain.Cryptocurrency, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var result []*domain.Cryptocurrency

	for _, history := range c.coins {
		if len(history) > 0 {
			result = append(result, history[len(history)-1])
		}
	}

	return result, nil
}

func (c *CryptoRepo) Get(ctx context.Context, symbol string) (*domain.Cryptocurrency, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	coin, exitsts := c.coins[symbol]
	if !exitsts {
		return nil, domain.ErrCryptoNotFound
	}

	return coin[len(coin)-1], nil
}

func (c *CryptoRepo) Add(ctx context.Context, crypto *domain.Cryptocurrency) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.coins[crypto.Symbol] == nil {
		c.coins[crypto.Symbol] = make([]*domain.Cryptocurrency, 0, 100)
	}

	c.coins[crypto.Symbol] = append(c.coins[crypto.Symbol], crypto)

	if len(c.coins[crypto.Symbol]) > 100 {
		c.coins[crypto.Symbol] = c.coins[crypto.Symbol][1:]
	}
}

func (c *CryptoRepo) GetHistory(ctx context.Context, symbol string) ([]*domain.Cryptocurrency, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	history, exists := c.coins[symbol]

	if !exists || len(history) == 0 {
		return nil, domain.ErrCryptoNotFound
	}

	result := make([]*domain.Cryptocurrency, len(history))
	copy(result, history)

	return result, nil
}

func (c *CryptoRepo) Delete(ctx context.Context, symbol string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, exists := c.coins[symbol]
	if !exists {
		return domain.ErrCryptoNotFound
	}

	delete(c.coins, symbol)

	return nil
}
