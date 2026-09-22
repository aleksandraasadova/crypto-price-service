package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/aleksandraasadova/crypto-price-service/internal/domain"
)

type Client struct {
	httpClient *http.Client
	baseURL    string

	mu        sync.RWMutex
	idCache   map[string]string // 'btc' -> 'bitcoin'
	nameCache map[string]string // 'btc' -> 'Bitcoin'
}

func NewClient() *Client {
	c := &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    "https://api.coingecko.com/api/v3",

		idCache:   make(map[string]string),
		nameCache: make(map[string]string),
	}
	/*
		if err := c.LoadCoinList(); err != nil {
			slog.Error("error with loading the coin list", "err", err)
		} else {
			slog.Info("coin list loaded succesfully")
		}*/

	return c
}

func (c *Client) GetCryptoCurrency(ctx context.Context, symbol string) (*domain.Cryptocurrency, error) {
	symbolLower := strings.ToLower(symbol)

	c.mu.RLock()
	coinID, idExists := c.idCache[symbolLower]
	coinName, nameExists := c.nameCache[symbolLower]
	c.mu.RUnlock()

	if !idExists || !nameExists {
		searchURL := fmt.Sprintf("%s/search?query=%s", c.baseURL, symbolLower)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL, nil)
		if err != nil {
			return nil, fmt.Errorf("search request error: %w", err)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("search request error: %w", err)
		}
		defer resp.Body.Close()

		var searchResult struct {
			Coins []struct {
				ID     string `json:"id"`
				Symbol string `json:"symbol"`
				Name   string `json:"name"`
			} `json:"coins"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
			return nil, fmt.Errorf("couldn't read search JSON: %w", err)
		}

		found := false
		for _, coin := range searchResult.Coins {
			if strings.EqualFold(coin.Symbol, symbol) {
				coinID = coin.ID
				coinName = coin.Name
				found = true
				break
			}
		}

		if !found {
			return nil, domain.ErrUnknownSymbol
		}

		c.mu.Lock()
		c.idCache[symbolLower] = coinID
		c.nameCache[symbolLower] = coinName
		c.mu.Unlock()
	}

	priceURL := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=usd", c.baseURL, coinID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, priceURL, nil)
	if err != nil {
		return nil, fmt.Errorf("price request error: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("price request error: %w", err)
	}
	defer resp.Body.Close()

	var priceResult map[string]struct {
		USD float64 `json:"usd"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&priceResult); err != nil {
		return nil, fmt.Errorf("couldn't read price from JSON: %w", err)
	}

	price := priceResult[coinID].USD

	return &domain.Cryptocurrency{
		Symbol:       symbol,
		Name:         coinName,
		CurrentPrice: price,
		LastUpdated:  time.Now(),
	}, nil
}

/*
// ID + Name
func (c *Client) LoadCoinList() error {
	url := fmt.Sprintf("%s/coins/list", c.baseURL)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	var coins []struct {
		ID     string `json:"id"`
		Symbol string `json:"symbol"`
		Name   string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&coins); err != nil {
		return fmt.Errorf("failed to read JSON: %w", err)
	}

	c.mu.Lock()
	for _, coin := range coins {
		symbolLower := strings.ToLower(coin.Symbol)

		if _, exists := c.idCache[symbolLower]; !exists {
			c.idCache[symbolLower] = coin.ID
			c.nameCache[symbolLower] = coin.Name
		}
	}
	c.mu.Unlock()

	return nil
}
*/
