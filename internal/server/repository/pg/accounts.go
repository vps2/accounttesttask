package pg

import (
	"context"

	"github.com/vps2/accounttesttask/internal/server/model"
	"github.com/vps2/accounttesttask/internal/server/repository"

	"github.com/go-pg/pg/v10"
)

type AccountsRepo struct {
	db *pg.DB
}

func NewAccountsRepo(db *pg.DB) *AccountsRepo {
	return &AccountsRepo{
		db: db,
	}
}

func (repo *AccountsRepo) GetById(ctx context.Context, id int32) (*model.Account, error) {
	account := &model.DBAccount{}
	err := repo.db.ModelContext(ctx, account).
		Where("id = ?", id).
		Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, repository.ErrAccountNotFound
		}

		return nil, err
	}

	return account.ToAccount(), nil
}

func (repo *AccountsRepo) Create(ctx context.Context, account *model.Account) (*model.Account, error) {
	dbAccount := account.ToDBAccount()
	_, err := repo.db.ModelContext(ctx, dbAccount).
		Insert()
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (repo *AccountsRepo) Update(ctx context.Context, account *model.Account) (*model.Account, error) {
	dbAccount := account.ToDBAccount()
	_, err := repo.db.ModelContext(ctx, dbAccount).
		WherePK().
		Update()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, repository.ErrAccountNotFound
		}

		return nil, err
	}

	return account, nil
}
