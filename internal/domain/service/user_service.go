package service

import (
	"context"
	"errors"
	"time"

	"github.com/elaurentium/listener-net/internal/domain/entities"
	"github.com/elaurentium/listener-net/internal/domain/repository"
	"github.com/google/uuid"
)

type UserService struct {
	userRepo repository.UserRepository
}


func NewUserRepository(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) Register(ctx context.Context, ip string, name string, dispositive string) (*entities.User, error) {
	exist, err := s.userRepo.CheckIPWasRegistred(ctx, ip)

	if err != nil {
		return nil, err
	}

	if exist {
		return nil, errors.New("ip already registred")
	}

	user := &entities.User{
		ID: uuid.New(),
		IP: ip,
		Name: name,
		LastSeen: time.Now(),
		Dispositive: dispositive,
	}

	err = s.userRepo.Create(ctx, user)

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