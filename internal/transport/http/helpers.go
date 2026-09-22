package transport

import (
	"encoding/json"
	"net/http"

	"github.com/aleksandraasadova/crypto-price-service/internal/domain"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func toCryptoDTO(crypto *domain.Cryptocurrency) *CryptoDTO {
	if crypto == nil {
		return nil
	}
	return &CryptoDTO{
		Symbol:       crypto.Symbol,
		Name:         crypto.Name,
		CurrentPrice: crypto.CurrentPrice,
		LastUpdated:  crypto.LastUpdated,
	}
}
