package storage

import "context"

type TransactorInterface interface {
	MakeTransaction(context.Context, func(ctx context.Context) error) error
}
