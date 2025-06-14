package repository

import (
	"context"

	"github.com/elaurentium/listener-net/internal/domain/entities"
	"github.com/google/uuid"
)


type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetByIP(ctx context.Context, ip string) (*entities.User, error)
	GetByDispositive(ctx context.Context, dispositive string) (*entities.User, error)
	CheckIPWasRegistered(ctx context.Context, ip string) (bool, error)
	Create(ctx context.Context, user *entities.User) error
}