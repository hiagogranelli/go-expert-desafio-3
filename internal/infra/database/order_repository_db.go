package database

import (
	"context"
	"database/sql"

	"desafio-cleanarchitecture/internal/domain/entity"
	"desafio-cleanarchitecture/internal/domain/repository"

	_ "github.com/lib/pq" // Driver PostgreSQL
	// ou _ "github.com/jackc/pgx/v5/stdlib"
)

type OrderRepositoryDb struct {
	Db *sql.DB // Ou *sql.Tx se usando Unit of Work
}

func NewOrderRepositoryDb(db *sql.DB) repository.OrderRepository {
	return &OrderRepositoryDb{Db: db}
}

func (r *OrderRepositoryDb) Save(ctx context.Context, order *entity.Order) error {
	// Usar transação se Db for *sql.Tx
	stmt, err := r.Db.PrepareContext(ctx, `
        INSERT INTO orders (id, price, tax, final_price, created_at, updated_at)
        VALUES ($1, $2, $3, $4, NOW(), NOW())
        ON CONFLICT (id) DO UPDATE SET
            price = $2,
            tax = $3,
            final_price = $4,
            updated_at = NOW()`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, order.ID, order.Price, order.Tax, order.FinalPrice)
	return err
}

func (r *OrderRepositoryDb) FindAll(ctx context.Context) ([]*entity.Order, error) {
	rows, err := r.Db.QueryContext(ctx, "SELECT id, price, tax, final_price, created_at, updated_at FROM orders ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*entity.Order
	for rows.Next() {
		var order entity.Order
		err := rows.Scan(
			&order.ID,
			&order.Price,
			&order.Tax,
			&order.FinalPrice,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			// Considerar logar o erro aqui
			continue // Ou retornar erro? Depende da política
		}
		orders = append(orders, &order)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}
