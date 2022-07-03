package repository

import (
	"context"
	"fmt"

	"github.com/maypok86/wb-l0/internal/entity"
	"github.com/maypok86/wb-l0/pkg/postgres"
)

type OrderPostgresRepository struct {
	db *postgres.Postgres
}

func NewOrderPostgresRepository(db *postgres.Postgres) OrderPostgresRepository {
	return OrderPostgresRepository{db: db}
}

func (opr OrderPostgresRepository) CreateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	tx, err := opr.db.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("can not begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	sql, args, err := opr.db.Builder.Insert("deliveries").Columns(
		"name",
		"phone",
		"zip",
		"city",
		"address",
		"region",
		"email",
	).Suffix("RETURNING id").Values(
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("can not build insert delivery query: %w", err)
	}

	var deliveryID int
	if err := tx.QueryRow(ctx, sql, args...).Scan(&deliveryID); err != nil {
		return nil, fmt.Errorf("")
	}

	sql, args, err = opr.db.Builder.Insert("payments").Columns(
		"transaction",
		"request_id",
		"currency",
		"provider",
		"amount",
		"payment_dt",
		"bank",
		"delivery_cost",
		"goods_total",
		"custom_fee",
	).Values(
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDt,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("can not build insert payment query: %w", err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("can not insert payment: %w", err)
	}

	sql, args, err = opr.db.Builder.Insert("orders").Columns(
		"order_uid",
		"track_number",
		"entry",
		"delivery_id",
		"locale",
		"internal_signature",
		"customer_id",
		"delivery_service",
		"shardkey",
		"sm_id",
		"date_created",
		"oof_shard",
	).Values(
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		deliveryID,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.Shardkey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("can not build insert order query: %w", err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("can not insert order: %w", err)
	}

	for _, item := range order.Items {
		sql, args, err = opr.db.Builder.Insert("items").Columns(
			"chrt_id",
			"track_number",
			"price",
			"rid",
			"name",
			"sale",
			"size",
			"total_price",
			"nm_id",
			"brand",
			"status",
		).Values(
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
		).ToSql()
		if err != nil {
			return nil, fmt.Errorf("can not build insert item query: %w", err)
		}

		_, err = tx.Exec(ctx, sql, args...)
		if err != nil {
			return nil, fmt.Errorf("can not insert item: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("can not commit transaction: %w", err)
	}

	return order, nil
}

func (opr OrderPostgresRepository) getItemsByTrackNumber(ctx context.Context, trackNumber string) ([]entity.Item, error) {
	sql, args, err := opr.db.Builder.Select(
		"chrt_id",
		"track_number",
		"price",
		"rid",
		"name",
		"sale",
		"size",
		"total_price",
		"nm_id",
		"brand",
		"status",
	).From("items").Where("track_number = ?", trackNumber).ToSql()
	if err != nil {
		return nil, fmt.Errorf("can not build select items query: %w", err)
	}

	rows, err := opr.db.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("can not select items: %w", err)
	}
	defer rows.Close()

	const defaultItemsCapacity = 64
	items := make([]entity.Item, 0, defaultItemsCapacity)

	for rows.Next() {
		var item entity.Item
		if err := rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		); err != nil {
			return nil, fmt.Errorf("can not scan item: %w", err)
		}

		items = append(items, item)
	}

	return items, nil
}

func (opr OrderPostgresRepository) GetOrderByID(ctx context.Context, orderUID string) (*entity.Order, error) {
	order := &entity.Order{}
	sql, args, err := opr.db.Builder.Select(
		"deliveries.name",
		"deliveries.phone",
		"deliveries.zip",
		"deliveries.city",
		"deliveries.address",
		"deliveries.region",
		"deliveries.email",

		"payments.transaction",
		"payments.request_id",
		"payments.currency",
		"payments.provider",
		"payments.amount",
		"payments.payment_dt",
		"payments.bank",
		"payments.delivery_cost",
		"payments.goods_total",
		"payments.custom_fee",

		"orders.order_uid",
		"orders.track_number",
		"orders.entry",
		"orders.locale",
		"orders.internal_signature",
		"orders.customer_id",
		"orders.delivery_service",
		"orders.shardkey",
		"orders.sm_id",
		"orders.date_created",
		"orders.oof_shard",
	).From("orders").Join(
		"deliveries ON orders.delivery_id = deliveries.id",
	).Join(
		"payments ON orders.order_uid = payments.transaction",
	).Where("orders.order_uid = ?", orderUID).Limit(1).ToSql()
	if err != nil {
		return nil, fmt.Errorf("can not build select order by id query: %w", err)
	}

	if err := opr.db.Pool.QueryRow(ctx, sql, args...).Scan(
		&order.Delivery.Name,
		&order.Delivery.Phone,
		&order.Delivery.Zip,
		&order.Delivery.City,
		&order.Delivery.Address,
		&order.Delivery.Region,
		&order.Delivery.Email,

		&order.Payment.Transaction,
		&order.Payment.RequestID,
		&order.Payment.Currency,
		&order.Payment.Provider,
		&order.Payment.Amount,
		&order.Payment.PaymentDt,
		&order.Payment.Bank,
		&order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal,
		&order.Payment.CustomFee,

		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.Shardkey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard,
	); err != nil {
		return nil, fmt.Errorf("can not select order by id: %w", err)
	}

	order.Items, err = opr.getItemsByTrackNumber(ctx, order.TrackNumber)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (opr OrderPostgresRepository) GetAllOrders(ctx context.Context) ([]*entity.Order, error) {
	sql, args, err := opr.db.Builder.Select(
		"deliveries.name",
		"deliveries.phone",
		"deliveries.zip",
		"deliveries.city",
		"deliveries.address",
		"deliveries.region",
		"deliveries.email",

		"payments.transaction",
		"payments.request_id",
		"payments.currency",
		"payments.provider",
		"payments.amount",
		"payments.payment_dt",
		"payments.bank",
		"payments.delivery_cost",
		"payments.goods_total",
		"payments.custom_fee",

		"orders.order_uid",
		"orders.track_number",
		"orders.entry",
		"orders.locale",
		"orders.internal_signature",
		"orders.customer_id",
		"orders.delivery_service",
		"orders.shardkey",
		"orders.sm_id",
		"orders.date_created",
		"orders.oof_shard",
	).From("orders").Join(
		"deliveries ON orders.delivery_id = deliveries.id",
	).Join(
		"payments ON orders.order_uid = payments.transaction",
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("can not build select all orders query: %w", err)
	}

	rows, err := opr.db.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("can not select items: %w", err)
	}
	defer rows.Close()

	const defaultOrdersCapacity = 64
	orders := make([]*entity.Order, 0, defaultOrdersCapacity)

	for rows.Next() {
		order := &entity.Order{}

		if err := rows.Scan(
			&order.Delivery.Name,
			&order.Delivery.Phone,
			&order.Delivery.Zip,
			&order.Delivery.City,
			&order.Delivery.Address,
			&order.Delivery.Region,
			&order.Delivery.Email,

			&order.Payment.Transaction,
			&order.Payment.RequestID,
			&order.Payment.Currency,
			&order.Payment.Provider,
			&order.Payment.Amount,
			&order.Payment.PaymentDt,
			&order.Payment.Bank,
			&order.Payment.DeliveryCost,
			&order.Payment.GoodsTotal,
			&order.Payment.CustomFee,

			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.Shardkey,
			&order.SmID,
			&order.DateCreated,
			&order.OofShard,
		); err != nil {
			return nil, fmt.Errorf("can not scan order: %w", err)
		}

		order.Items, err = opr.getItemsByTrackNumber(ctx, order.TrackNumber)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}
