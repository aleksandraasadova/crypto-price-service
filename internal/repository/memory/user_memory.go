package memory

import (
	"context"
	"log/slog"
	"sync"

	"github.com/aleksandraasadova/crypto-price-service/internal/domain"
)

type UserRepo struct {
	mu    sync.RWMutex
	users map[string]*domain.User
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		users: make(map[string]*domain.User),
	}
}

func (m *UserRepo) Get(ctx context.Context, username string) (*domain.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if user, exists := m.users[username]; exists {
		return user, nil
	}
	slog.Debug("user doesnt exists - get repo")
	return nil, domain.ErrUserNotFound
}

func (m *UserRepo) Create(ctx context.Context, user *domain.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.users[user.Username]; exists {
		return domain.ErrUserAlreadyExists
	}

	m.users[user.Username] = user
	slog.Debug("user created")
	return nil
}

//accept interfaces return structs
