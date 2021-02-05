package service

import (
	"accounts/internal/server/model"
	"accounts/internal/server/repository"
	"accounts/pkg/cache"
	"context"
	"errors"
	"sync"
)

type AccountsSvc struct {
	mu    *sync.RWMutex
	repo  repository.Accounts
	cache cache.Cache
}

func NewAccountsSvc(repo repository.Accounts, cache cache.Cache) *AccountsSvc {
	return &AccountsSvc{
		mu:    &sync.RWMutex{},
		repo:  repo,
		cache: cache,
	}
}

func (svc *AccountsSvc) GetAmount(ctx context.Context, id int32) (int64, error) {
	svc.mu.RLock()
	defer svc.mu.RUnlock()

	if val, ok := svc.cache.Get(id); ok {
		return val.(int64), nil
	}

	account, err := svc.repo.GetById(ctx, id)
	if err != nil {
		return 0, nil
	} else {
		return account.Balance, nil
	}
}

func (svc *AccountsSvc) AddAmount(ctx context.Context, id int32, amount int64) error {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	if account, err := svc.repo.GetById(ctx, id); err == nil { //нашли запись в хранилище
		newAmount := account.Balance + amount
		if newAmount < 0 {
			return errors.New("the balance is less than the withdrawal amount")
		}

		account.Balance = newAmount

		return svc.update(ctx, account)
	} else if err == repository.ErrAccountNotFound { //записи нет в хранилище
		if amount <= 0 {
			return errors.New("cannot create an account with a negative or zero balance")
		}

		account = &model.Account{Id: id, Balance: amount}

		return svc.create(ctx, account)
	} else {
		return err
	}
}

func (svc *AccountsSvc) update(ctx context.Context, account *model.Account) error {
	_, err := svc.repo.Update(ctx, account)
	if err != nil {
		return err
	}

	svc.cache.Set(account.Id, account.Balance)

	return nil
}

func (svc *AccountsSvc) create(ctx context.Context, account *model.Account) error {
	_, err := svc.repo.Create(ctx, account)
	if err != nil {
		return err
	}

	svc.cache.Set(account.Id, account.Balance)

	return nil
}
