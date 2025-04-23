package repository

import (
	"context"
	"database/sql"

	"desafio-cleanarchitecture/internal/domain/entity"
)

type OrderRepository interface {
	Save(ctx context.Context, order *entity.Order) error
	FindAll(ctx context.Context) ([]*entity.Order, error)
}

type Uow interface {
	Do(ctx context.Context, fn func(uow *Uow) error) error
	Register(name string, fc RepositoryFactory)
	GetRepository(ctx context.Context, name string) (interface{}, error)
	CommitOrRollback() error
	Rollback() error
	UnRegister(name string)
}

type RepositoryFactory func(tx *sql.Tx) interface{}
