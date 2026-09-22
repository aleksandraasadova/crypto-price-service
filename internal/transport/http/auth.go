package transport

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/aleksandraasadova/crypto-price-service/internal/domain"
)

type AuthService interface {
	Register(ctx context.Context, username, password string) (string, error)
	Login(ctx context.Context, username, password string) (string, error)
}

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid json body"}) // 400
		return
	}

	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "username and password required"})
		return
	}

	token, err := a.authService.Register(r.Context(), req.Username, req.Password)
	if err != nil { // 400 or 409 // 409 - пользователь уже существует
		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid credentials"})
		case errors.Is(err, domain.ErrUserAlreadyExists):
			WriteJSON(w, http.StatusConflict, ErrorResponse{Error: "username already taken"})
		default:
			log.Printf("Register handler error: %v", err)
			WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		}
		return
	}

	WriteJSON(w, http.StatusCreated, RegisterUserResponse{Token: token})
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid json body"}) // 400
		return
	}

	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "username and password required"})
		return
	}

	token, err := a.authService.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid credentials"})
			return
		}
	}

	WriteJSON(w, http.StatusOK, LoginUserResponse{Token: token})
}
