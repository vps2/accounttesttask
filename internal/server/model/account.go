package model

type Account struct {
	Id      int32
	Balance int64
}

func (account *Account) ToDBAccount() *DBAccount {
	return &DBAccount{
		Id:      account.Id,
		Balance: account.Balance,
	}
}

// DBAccount is a Postgres user
type DBAccount struct {
	tableName struct{} `pg:"accounts"`
	Id        int32    `pg:",notnull,pk"`
	Balance   int64    `pg:",use_zero,notnull"`
}

func (dbAccount *DBAccount) ToAccount() *Account {
	return &Account{
		Id:      dbAccount.Id,
		Balance: dbAccount.Balance,
	}
}
