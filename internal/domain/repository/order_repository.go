package repository

import (
	"context"
	"database/sql" // Ou interface específica se usar Unit of Work

	"desafio-cleanarchitecture/internal/domain/entity"
)

type OrderRepository interface {
	Save(ctx context.Context, order *entity.Order) error
	FindAll(ctx context.Context) ([]*entity.Order, error)
	// FindByID(ctx context.Context, id string) (*entity.Order, error) // Opcional
}

// Interface para Unit of Work (Opcional mas recomendado)
type Uow interface {
	Do(ctx context.Context, fn func(uow *Uow) error) error
	Register(name string, fc RepositoryFactory)
	GetRepository(ctx context.Context, name string) (interface{}, error)
	CommitOrRollback() error
	Rollback() error // Adicionado para clareza
	UnRegister(name string)
}

type RepositoryFactory func(tx *sql.Tx) interface{}
