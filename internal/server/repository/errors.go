package repository

import "errors"

//Ошибки, которые могут возвратить экземпляры repository.Accounts
var (
	ErrAccountNotFound      = errors.New("account not found")
	ErrAccountAlreadyExists = errors.New("account already exists")
)
