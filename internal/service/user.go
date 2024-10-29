package service

import (
	"assignment-pe/internal/errs"
	"assignment-pe/internal/model"
	"assignment-pe/internal/repo"
	"context"
)

// nolint: lll
type UserService interface {
	GetUserByAddress(ctx context.Context, address string) (*model.User, error)
}

type userService struct {
	userRepo repo.UserRepo
}

func NewUserService(
	userRepo repo.UserRepo,
) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (srv *userService) GetUserByAddress(
	ctx context.Context,
	address string,
) (*model.User, error) {
	user, err := srv.userRepo.GetByAddress(ctx, address)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errs.ErrNotFound.ReMsgf("user %v not found", address)
	}

	return user, nil
}
