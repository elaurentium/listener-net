package service

import (
	"context"
	"errors"
	"time"

	"github.com/elaurentium/listener-net/internal/domain/entities"
	"github.com/elaurentium/listener-net/internal/domain/repository"
	"github.com/google/uuid"
)


type BlacklistService struct {
	blacklistRepo repository.BlacklistRepository
}

func NewBlacklistService(blacklistRepo repository.BlacklistRepository) *BlacklistService {
	return &BlacklistService{
		blacklistRepo: blacklistRepo,
	}
}

type BlacklistRequest struct {
	IP string `json:"ip"`
	Name string `json:"name"`
	Dispositive string `json:"dispositive"`
}

func (r *BlacklistService) Register(ctx context.Context, ip string, name string, dispositive string) (*entities.Blacklist, error) {
	exist, err := r.blacklistRepo.CheckIPWasRegistred(ctx, ip)

	if err != nil {
		return nil, err
	}

	if exist {
		return nil, errors.New("ip already registred on blacklist")
	}

	blacklist := &entities.Blacklist{
		ID: uuid.New(),
		IP: ip,
		Name: name,
		LastSeen: time.Now(),
		Dispositive: dispositive,
	}

	err = r.blacklistRepo.Create(ctx, blacklist)

	if err != nil {
		return nil, err
	}

	return blacklist, nil
}