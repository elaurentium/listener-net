package db

import (
	"context"

	"github.com/elaurentium/listener-net/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepository interface {
	GetById(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetByIP(ctx context.Context, ip string) (*entities.User, error)
	Create(ctx context.Context, user *entities.User) error
}


type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	//return &userRepository{
	//	pool: pool,
	//}
}