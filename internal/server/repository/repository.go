package repository

import (
	"accounts/internal/server/model"
	"context"
)

//go:generate mockery --dir . --name Accounts --filename accounts.go --structname AccountsRepo --output ./mocks
type Accounts interface {
	GetById(context.Context, int32) (*model.Account, error)
	Create(context.Context, *model.Account) (*model.Account, error)
	Update(context.Context, *model.Account) (*model.Account, error)
}
