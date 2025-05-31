package service

import (
	"context"

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

func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	return s.userRepo.GetByIP(ctx, id.String())
}