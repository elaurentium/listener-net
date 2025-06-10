package repository

import (
	"context"

	"github.com/elaurentium/listener-net/internal/domain/entities"
	"github.com/google/uuid"
)

type BlacklistRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Blacklist, error)
	GetByIP(ctx context.Context, ip string) (*entities.Blacklist, error)
	CheckIPWasRegistered(ctx context.Context, ip string) (bool, error)
	Create(ctx context.Context, user *entities.Blacklist) error
}