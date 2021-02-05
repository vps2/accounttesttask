package service

import "context"

type AccountsService interface {
	GetAmount(ctx context.Context, id int32) (int64, error)
	AddAmount(ctx context.Context, id int32, amount int64) error
}

type StatisticsService interface {
	IncReadOperations()
	IncWriteOperations()
	Reset()
	TotalReadOperations() int64
	TotalWriteOperations() int64
	ReadOperationsPerSecond() int64
	WriteOperationsPerSecond() int64
}
