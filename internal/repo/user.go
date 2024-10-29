package repo

import (
	"assignment-pe/internal/cx"
	"assignment-pe/internal/errs"
	"assignment-pe/internal/model"
	"context"
	"database/sql"
)

type UserRepo interface {
	GetByAddress(ctx context.Context, address string) (*model.User, error)
}

type userRepo struct {
}

func NewUserRepo() UserRepo {
	return &userRepo{}
}

type dbUser struct {
	ID      string `db:"id"`
	Address string `db:"address"`
}

func (d dbUser) toEntity() model.User {
	return model.User(d)
}

func (repo *userRepo) GetByAddress(
	ctx context.Context,
	address string,
) (*model.User, error) {
	tx := cx.GetTx(ctx)

	stms := `
		SELECT id, address
		FROM users
		WHERE address = ?
	`
	var user dbUser
	err := tx.GetContext(ctx, &user, tx.Rebind(stms), address)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errs.ErrInternal.Rewrap(err)
	}

	userEntity := user.toEntity()
	return &userEntity, nil
}
