package service

import (
	"accounts/internal/server/model"
	"accounts/internal/server/repository"
	rmocks "accounts/internal/server/repository/mocks"
	cmocks "accounts/pkg/cache/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"gotest.tools/assert"
)

func TestAccountsSvc_GetAmount(t *testing.T) {
	input := &model.Account{
		Id:      1,
		Balance: 300,
	}

	tests := []struct {
		name         string
		expectations func(accountsRepo *rmocks.AccountsRepo, cache *cmocks.Cache)
		input        int32
		want         int64
	}{
		{
			name: "value in cache",
			expectations: func(accountsRepo *rmocks.AccountsRepo, cache *cmocks.Cache) {
				cache.On("Get", input.Id).Return(input.Balance, true)
			},
			input: input.Id,
			want:  input.Balance,
		},
		{
			name: "in db",
			expectations: func(accountsRepo *rmocks.AccountsRepo, cache *cmocks.Cache) {
				cache.On("Get", mock.AnythingOfType("int32")).Return(nil, false)
				accountsRepo.On("GetById", context.Background(), input.Id).Return(input, nil)
			},
			input: input.Id,
			want:  input.Balance,
		},
		{
			name: "not in db",
			expectations: func(accountsRepo *rmocks.AccountsRepo, cache *cmocks.Cache) {
				cache.On("Get", mock.AnythingOfType("int32")).Return(nil, false)
				accountsRepo.On("GetById", context.Background(), input.Id).Return(nil, errors.New("account not found"))
			},
			input: input.Id,
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			accountsRepo := &rmocks.AccountsRepo{}
			cache := &cmocks.Cache{}
			svc := NewAccountsSvc(accountsRepo, cache)
			tt.expectations(accountsRepo, cache)

			if got, _ := svc.GetAmount(ctx, tt.input); got != tt.want {
				t.Errorf("AccountsSvc.GetAmount(%v) = %v, want %v", tt.input, got, tt.want)
			}

			accountsRepo.AssertExpectations(t)
			cache.AssertExpectations(t)
		})
	}
}

func TestAccountsSvc_AddAmount(t *testing.T) {
	tests := []struct {
		name         string
		expectations func(accountsRepo *rmocks.AccountsRepo, cache *cmocks.Cache)
		input        *model.Account
		err          error
	}{
		{
			name: "new account with positive balance",
			expectations: func(accountsRepo *rmocks.AccountsRepo, cache *cmocks.Cache) {
				in := &model.Account{
					Id:      1,
					Balance: 300,
				}
				cache.On("Set", mock.AnythingOfType("int32"), mock.AnythingOfType("int64")).Return(true)
				accountsRepo.On("GetById", context.Background(), in.Id).Return(nil, repository.ErrAccountNotFound)
				accountsRepo.On("Create", context.Background(), in).Return(in, nil)
			},
			input: &model.Account{
				Id:      1,
				Balance: 300,
			},
		},
		{
			name: "new account with negative balance",
			expectations: func(accountsRepo *rmocks.AccountsRepo, cache *cmocks.Cache) {
				id := int32(1)
				accountsRepo.On("GetById", context.Background(), id).Return(nil, repository.ErrAccountNotFound)
			},
			input: &model.Account{
				Id:      1,
				Balance: -300,
			},
			err: errors.New("cannot create an account with a negative or zero balance"),
		},
		{
			name: "top up your account balance.",
			expectations: func(accountsRepo *rmocks.AccountsRepo, cache *cmocks.Cache) {
				existsAccount := &model.Account{
					Id:      1,
					Balance: 300,
				}
				cache.On("Set", mock.AnythingOfType("int32"), mock.AnythingOfType("int64")).Return(true)
				accountsRepo.On("GetById", context.Background(), existsAccount.Id).Return(existsAccount, nil)
				accountsRepo.On("Update", context.Background(), existsAccount).Return(existsAccount, nil)
			},
			input: &model.Account{
				Id:      1,
				Balance: 300,
			},
		},
		{
			name: "withdraw an amount greater than the balance",
			expectations: func(accountsRepo *rmocks.AccountsRepo, cache *cmocks.Cache) {
				existsAccount := &model.Account{
					Id:      1,
					Balance: 300,
				}
				accountsRepo.On("GetById", context.Background(), existsAccount.Id).Return(existsAccount, nil)
			},
			input: &model.Account{
				Id:      1,
				Balance: -400,
			},
			err: errors.New("the balance is less than the withdrawal amount"),
		},
		{
			name: "withdraw the entire amount from the account",
			expectations: func(accountsRepo *rmocks.AccountsRepo, cache *cmocks.Cache) {
				existsAccount := &model.Account{
					Id:      1,
					Balance: 300,
				}
				cache.On("Set", mock.AnythingOfType("int32"), mock.AnythingOfType("int64")).Return(true)
				accountsRepo.On("GetById", context.Background(), existsAccount.Id).Return(existsAccount, nil)
				accountsRepo.On("Update", context.Background(), existsAccount).Return(existsAccount, nil)
			},
			input: &model.Account{
				Id:      1,
				Balance: -300,
			},
		},
		{
			name: "repo error on create",
			expectations: func(accountsRepo *rmocks.AccountsRepo, cache *cmocks.Cache) {
				in := &model.Account{
					Id:      1,
					Balance: 300,
				}
				accountsRepo.On("GetById", context.Background(), in.Id).Return(nil, repository.ErrAccountNotFound)
				accountsRepo.On("Create", context.Background(), in).Return(nil, errors.New("some error"))
			},
			input: &model.Account{
				Id:      1,
				Balance: 300,
			},
			err: errors.New("some error"),
		},
		{
			name: "repo error on update",
			expectations: func(accountsRepo *rmocks.AccountsRepo, cache *cmocks.Cache) {
				existingAccount := &model.Account{
					Id:      1,
					Balance: 300,
				}
				accountsRepo.On("GetById", context.Background(), existingAccount.Id).Return(existingAccount, nil)
				accountsRepo.On("Update", context.Background(), existingAccount).Return(nil, errors.New("some error"))
			},
			input: &model.Account{
				Id:      1,
				Balance: 300,
			},
			err: errors.New("some error"),
		},
		{
			name: "repo error on get",
			expectations: func(accountsRepo *rmocks.AccountsRepo, cache *cmocks.Cache) {
				in := &model.Account{
					Id:      1,
					Balance: 300,
				}
				accountsRepo.On("GetById", context.Background(), in.Id).Return(nil, errors.New("some error"))
			},
			input: &model.Account{
				Id:      1,
				Balance: 300,
			},
			err: errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			accountsRepo := &rmocks.AccountsRepo{}
			cache := &cmocks.Cache{}
			svc := NewAccountsSvc(accountsRepo, cache)
			tt.expectations(accountsRepo, cache)

			err := svc.AddAmount(ctx, tt.input.Id, tt.input.Balance)
			if err != nil {
				if tt.err != nil {
					assert.Equal(t, tt.err.Error(), err.Error())
				} else {
					t.Errorf("AccountsSvc.AddAmount(%v, %v):  expected no error, found: %s",
						tt.input.Id,
						tt.input.Balance,
						err.Error())
				}
			}

			accountsRepo.AssertExpectations(t)
			cache.AssertExpectations(t)
		})
	}
}
