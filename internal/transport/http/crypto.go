package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aleksandraasadova/crypto-price-service/internal/domain"
)

type CryptoService interface {
	GetCryptos(ctx context.Context) ([]*domain.Cryptocurrency, error)
	AddCrypto(ctx context.Context, symbol string) (*domain.Cryptocurrency, error)
	GetCrypto(ctx context.Context, symbol string) (*domain.Cryptocurrency, error)
	RefreshCrypto(ctx context.Context, symbol string) (*domain.Cryptocurrency, error)
	GetHistory(ctx context.Context, symbol string) ([]*domain.Cryptocurrency, error)
	GetStats(ctx context.Context, symbol string) (*domain.CryptoStats, error)
	DeleteCrypto(ctx context.Context, symbol string) error
}

type CryptoHandler struct {
	cryptoService CryptoService
}

func NewCryptoHandler(cryptoService CryptoService) *CryptoHandler {
	return &CryptoHandler{
		cryptoService: cryptoService,
	}
}

// GET /crypto - получить список всех криптовалют
func (c *CryptoHandler) GetCryptos(w http.ResponseWriter, r *http.Request) {
	cryptos, err := c.cryptoService.GetCryptos(r.Context())
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "server error"})
		return
	}

	cryptosDTO := make([]*CryptoDTO, 0, len(cryptos))

	for _, crypto := range cryptos {
		cryptosDTO = append(cryptosDTO, toCryptoDTO(crypto))
	}

	WriteJSON(w, http.StatusOK, GetCryptosResponse{Cryptos: cryptosDTO})
}

// POST /crypto - добавить новую криптовалюту для отслеживания
func (c *CryptoHandler) AddCrypto(w http.ResponseWriter, r *http.Request) {
	var req AddCryptoRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid json body"})
		return
	}

	if req.Symbol == "" {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "symbol required"})
		return
	}

	crypto, err := c.cryptoService.AddCrypto(r.Context(), req.Symbol)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrCryptoExists):
			WriteJSON(w, http.StatusConflict, ErrorResponse{Error: "crypto currency already exists"})
			return
		case errors.Is(err, domain.ErrUnknownSymbol):
			WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "unknown crypto symbol"})
			return
		default:
			WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "server error"})
			return
		}
	}

	WriteJSON(w, http.StatusCreated, AddCryptoResponse{Crypto: toCryptoDTO(crypto)})
}

// GET /crypto/{symbol} - получить информацию о конкретной криптовалюте
func (c *CryptoHandler) GetCrypto(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")

	crypto, err := c.cryptoService.GetCrypto(r.Context(), symbol)
	if err != nil {
		if errors.Is(err, domain.ErrCryptoNotFound) {
			WriteJSON(w, http.StatusNotFound, ErrorResponse{Error: "no crypto currency with provided symbol found"})
			return
		}

		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "server error"})
		return
	}

	WriteJSON(w, http.StatusOK, toCryptoDTO(crypto))
}

// PUT /crypto/{symbol}/refresh - принудительно обновить цену криптовалюты
func (c *CryptoHandler) RefreshCrypto(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")

	crypto, err := c.cryptoService.RefreshCrypto(r.Context(), symbol)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrCryptoNotFound):
			WriteJSON(w, http.StatusNotFound, ErrorResponse{Error: "no crypto currency with provided symbol found"})
			return
		default:
			WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "server error"})
			return
		}
	}

	WriteJSON(w, http.StatusOK, RefreshCryptoResponse{Crypto: toCryptoDTO(crypto)})
}

// GET /crypto/{symbol}/history - получить историю цен криптовалюты
func (c *CryptoHandler) GetCryptoHistory(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")

	history, err := c.cryptoService.GetHistory(r.Context(), symbol)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrCryptoNotFound):
			WriteJSON(w, http.StatusNotFound, ErrorResponse{Error: "cryptocurrency not found"})
			return
		default:
			WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "server error"})
			return
		}
	}

	historyDTO := make([]*CryptoHistoryDTO, len(history))
	for i, v := range history {
		historyDTO[i] = &CryptoHistoryDTO{
			Price:     v.CurrentPrice,
			Timestamp: v.LastUpdated,
		}
	}

	WriteJSON(w, http.StatusOK, GetHistory{
		Symbol:  symbol,
		History: historyDTO,
	})
}

// GET /crypto/{symbol}/stats - получить статистику по ценам криптовалюты
func (c *CryptoHandler) GetCryptoStats(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")

	stats, err := c.cryptoService.GetStats(r.Context(), symbol)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrCryptoNotFound):
			WriteJSON(w, http.StatusNotFound, ErrorResponse{Error: "cryptocurrency not found"})
			return
		default:
			WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "server error"})
			return
		}
	}

	response := GetCryptoStatsResponse{
		Symbol:       symbol,
		CurrentPrice: stats.CurrentPrice,
		Stats: &StatsDTO{
			MinPrice:       stats.MinPrice,
			MaxPrice:       stats.MaxPrice,
			AvgPrice:       stats.AvgPrice,
			PriceChange:    stats.PriceChange,
			PriceChangePct: stats.PriceChangePct,
			RecordsCount:   stats.RecordsCount,
		},
	}

	WriteJSON(w, http.StatusOK, response)
}

func (c *CryptoHandler) DeleteCrypto(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")

	err := c.cryptoService.DeleteCrypto(r.Context(), symbol)
	if err != nil {
		if errors.Is(err, domain.ErrCryptoNotFound) {
			WriteJSON(w, http.StatusNotFound, ErrorResponse{Error: "cryptocurrency not found"})
			return
		}
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "server error"})
		return
	}

	WriteJSON(w, http.StatusOK, map[string]any{})

}
