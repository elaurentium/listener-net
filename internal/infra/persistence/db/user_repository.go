package db

import (
	"context"
	"fmt"

	"github.com/elaurentium/listener-net/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetByIP(ctx context.Context, ip string) (*entities.User, error)
	GetByDispositive(ctx context.Context, dispositive string) (*entities.User, error)
	Create(ctx context.Context, user *entities.User) error
	CheckIPWasRegistered(ctx context.Context, ip string) (bool, error)
}

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{
		pool: pool,
	}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO USERS (ID, IP, NAME, LAST_SEEN, DISPOSITIVE)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (DISPOSITIVE) DO UPDATE SET
			IP = EXCLUDED.IP,
			LAST_SEEN = EXCLUDED.LAST_SEEN
	`
	_, err := r.pool.Exec(ctx, query, user.ID, user.IP, user.Name, user.LastSeen, user.Dispositive)

	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	user := &entities.User{}
	query := `
		SELECT * FROM USERS WHERE ID = $1
	`
	err := r.pool.QueryRow(ctx, query, id).Scan(&user.ID)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetByIP(ctx context.Context, ip string) (*entities.User, error) {
	user := &entities.User{}
	query := `
		SELECT * FROM USERS WHERE IP = $1
	`
	err := r.pool.QueryRow(ctx, query, ip).Scan(&user.IP)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetByDispositive(ctx context.Context, dispositive string) (*entities.User, error) {
	user := &entities.User{}
	query := `
		SELECT * FROM USERS WHERE DISPOSITIVE = $1
	`
	err := r.pool.QueryRow(ctx, query, dispositive).Scan(&user.Dispositive)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) CheckIPWasRegistered(ctx context.Context, ip string) (bool, error) {
	var exists bool

	err := r.pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM USERS WHERE IP = $1)", ip).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check if IP was registered: %w", err)
	}

	return exists, nil
}
