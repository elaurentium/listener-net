package service

import (
	"context"
	"time"

	"github.com/elaurentium/listener-net/internal/domain/entities"
	"github.com/elaurentium/listener-net/internal/domain/repository"
	"github.com/elaurentium/listener-net/internal/infra/web/auth"
	"github.com/google/uuid"
)

type UserService struct {
	userRepo repository.UserRepository
	auth auth.AuthService
}


func NewUserService(userRepo repository.UserRepository, auth auth.AuthService) *UserService {
	return &UserService{
		userRepo: userRepo,
		auth: auth,
	}
}

func (s *UserService) Register(ctx context.Context, ip string, name string, dispositive string) (*entities.User, error) {
	user := &entities.User{
		ID: uuid.New(),
		IP: ip,
		Name: name,
		LastSeen: time.Now(),
		Dispositive: dispositive,
	}

	err := s.userRepo.Create(ctx, user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) GetByIP(ctx context.Context, ip string) (*entities.User, error) {
	return s.userRepo.GetByIP(ctx, ip)
}

func (s *UserService) GetByDispositive(ctx context.Context, dispositive string) (*entities.User, error) {
	return s.userRepo.GetByDispositive(ctx, dispositive)
}