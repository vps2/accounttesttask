package inmem

import (
	"accounts/internal/server/model"
	"accounts/internal/server/repository"
	"context"
)

type AccountsRepo struct {
	m map[int32]*model.Account
}

func NewAccountsRepo() *AccountsRepo {
	return &AccountsRepo{
		m: make(map[int32]*model.Account),
	}
}

func (a *AccountsRepo) GetById(ctx context.Context, id int32) (*model.Account, error) {

	if account, ok := a.m[id]; ok {
		return account, nil
	}

	return nil, repository.ErrAccountNotFound
}

func (a *AccountsRepo) Create(ctx context.Context, account *model.Account) (*model.Account, error) {
	if _, ok := a.m[account.Id]; !ok {
		a.m[account.Id] = account

		return account, nil
	}

	return nil, repository.ErrAccountAlreadyExists
}

func (a *AccountsRepo) Update(ctx context.Context, account *model.Account) (*model.Account, error) {
	if acc, ok := a.m[account.Id]; ok {
		acc.Balance = account.Balance

		return acc, nil
	}

	return nil, repository.ErrAccountNotFound
}
