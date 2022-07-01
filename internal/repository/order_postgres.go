package repository

import (
	"context"
	"fmt"

	"github.com/maypok86/wb-l0/pkg/postgres"
)

type OrderPostgresRepository struct {
	db *postgres.Postgres
}

func NewOrderPostgresRepository(db *postgres.Postgres) OrderPostgresRepository {
	return OrderPostgresRepository{db: db}
}

func (opr OrderPostgresRepository) CreateOrder(ctx context.Context, data []byte) (int, error) {
	sql, args, err := opr.db.Builder.Insert("orders").Columns("data").Suffix("RETURNING id").Values(data).ToSql()
	if err != nil {
		return 0, fmt.Errorf("can not build sql for create order query: %w", err)
	}

	var id int
	if err := opr.db.Pool.QueryRow(ctx, sql, args...).Scan(&id); err != nil {
		return 0, fmt.Errorf("can not create order: %w", err)
	}
	return id, nil
}

func (opr OrderPostgresRepository) GetOrderByID(ctx context.Context, id int) ([]byte, error) {
	sql, args, err := opr.db.Builder.Select("data").From("orders").Where("id = ?", id).Limit(1).ToSql()
	if err != nil {
		return nil, fmt.Errorf("can not build sql for get order by id query: %w", err)
	}

	var data []byte
	if err := opr.db.Pool.QueryRow(ctx, sql, args...).Scan(&data); err != nil {
		return nil, fmt.Errorf("can not get order by id: %w", err)
	}
	return data, nil
}

func (opr OrderPostgresRepository) GetAllOrders(ctx context.Context) ([][]byte, error) {
	sql, _, err := opr.db.Builder.Select("data").From("orders").ToSql()
	if err != nil {
		return nil, fmt.Errorf("can not build sql for get all orders query: %w", err)
	}

	rows, err := opr.db.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("can not exec get all orders query: %w", err)
	}
	defer rows.Close()

	const defaultEntityCapacity = 64
	entities := make([][]byte, 0, defaultEntityCapacity)
	for rows.Next() {
		var e []byte
		if err := rows.Scan(&e); err != nil {
			return nil, fmt.Errorf("can not scan order: %w", err)
		}
		entities = append(entities, e)
	}
	return entities, nil
}
