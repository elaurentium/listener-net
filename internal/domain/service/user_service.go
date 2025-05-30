package service

import "github.com/elaurentium/listener-net/internal/domain/repository"

type UserService struct {
	userRepo repository.UserRepository
}


func NewUserRepository(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}