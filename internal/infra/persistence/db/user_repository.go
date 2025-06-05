package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/elaurentium/listener-net/internal/domain/entities"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetByIP(ctx context.Context, ip string) (*entities.User, error)
	GetByDispositive(ctx context.Context, dispositive string) (*entities.User, error)
	Create(ctx context.Context, user *entities.User) error
	CheckIPWasRegistred(ctx context.Context, ip string) (bool, error)
}

type userRepository struct {
	sql *sql.DB
}

func NewUserRepository(sql *sql.DB) UserRepository {
	return &userRepository{
		sql: sql,
	}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO USERS (ID, IP, NAME, EMAIL) 
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.sql.ExecContext(ctx, query, user.ID, user.IP, user.Name)

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
	err := r.sql.QueryRowContext(ctx, query, id).Scan(&user.ID)

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
	err := r.sql.QueryRowContext(ctx, query, ip).Scan(&user.IP)

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
	err := r.sql.QueryRowContext(ctx, query, dispositive).Scan(&user.Dispositive)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) CheckIPWasRegistred(ctx context.Context, ip string) (bool, error) {
	var exists bool

	err := r.sql.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM USERS WHERE IP = $1)", ip).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check if IP was registred: %w", err)
	}

	return exists, nil
}
