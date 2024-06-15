package repository

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"

	"L_0_WB/internal/entity"
	"L_0_WB/internal/repository/postgres"
)

type DataBaseOrderRepository struct {
	storage *postgres.Storage
}

func NewDataBaseOrderRepository(storage *postgres.Storage) *DataBaseOrderRepository {
	return &DataBaseOrderRepository{storage: storage}
}

func (d *DataBaseOrderRepository) GetAllOrders(ctx context.Context) ([]entity.Order, error) {
	var orders []entity.Order

	rows, err := d.storage.QueryContext(ctx, `SELECT 
		"order_uid", "track_number", "entry", "locale", "internal_signature", 
		"customer_id", "delivery_service", "shardkey", "sm_id", "date_created", "oof_shard" 
		FROM "Order"`)

	if err != nil {
		return nil, fmt.Errorf("error querying orders: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order entity.Order
		err := rows.Scan(
			&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale,
			&order.InternalSignature, &order.CustomerID, &order.DeliveryService,
			&order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard)

		if err != nil {
			return nil, fmt.Errorf("error scanning order: %w", err)
		}

		err = d.getDeliveryInfo(ctx, &order)
		if err != nil {
			return nil, err
		}

		err = d.getPaymentInfo(ctx, &order)
		if err != nil {
			return nil, err
		}

		err = d.getItemsInfo(ctx, &order)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through orders: %w", err)
	}

	return orders, nil
}

func (d *DataBaseOrderRepository) getDeliveryInfo(ctx context.Context, order *entity.Order) error {
	query := `SELECT "name", "phone", "zip", "city", "address", "region", "email" 
			  FROM "Delivery" WHERE "order_uid" = $1`
	return d.storage.QueryRowContext(ctx, query, order.OrderUID).Scan(
		&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip,
		&order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email)
}

func (d *DataBaseOrderRepository) getPaymentInfo(ctx context.Context, order *entity.Order) error {
	query := `SELECT "transaction", "request_id", "currency", "provider", "amount", 
					 "payment_dt", "bank", "delivery_cost", "goods_total", "custom_fee" 
			  FROM "Payment" WHERE "order_uid" = $1`
	return d.storage.QueryRowContext(ctx, query, order.OrderUID).Scan(
		&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency,
		&order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt,
		&order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee)
}

func (d *DataBaseOrderRepository) getItemsInfo(ctx context.Context, order *entity.Order) error {
	query := `SELECT "chrt_id", "track_number", "price", "rid", "name", "sale", 
					 "size", "total_price", "nm_id", "brand", "status" 
			  FROM "Item" WHERE "order_uid" = $1`
	rows, err := d.storage.QueryContext(ctx, query, order.OrderUID)
	if err != nil {
		return fmt.Errorf("error querying items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.Item
		err := rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name,
			&item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
			return fmt.Errorf("error scanning item: %w", err)
		}
		order.Items = append(order.Items, item)
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("error iterating through items: %w", err)
	}

	return nil
}

func (d *DataBaseOrderRepository) Get(ctx context.Context, orderUID string) (entity.Order, error) {
	var order entity.Order

	query := `SELECT * FROM "Order" WHERE "order_uid" = $1`

	if err := d.storage.QueryRowContext(ctx, query, orderUID).Scan(&order); err != nil {
		return order, fmt.Errorf("error scanning order: %w", err)
	}

	return order, nil
}

func (d *DataBaseOrderRepository) Set(ctx context.Context, order entity.Order) error {
	if err := d.setOrderInfo(ctx, order); err != nil {
		return fmt.Errorf("connot set order info: %w", err)
	}
	log.Info().Msg("successfully set order")

	if err := d.setDeliveryInfo(ctx, order); err != nil {
		return fmt.Errorf("connot set order info: %w", err)
	}
	log.Info().Msg("successfully delivery order")

	if err := d.setPaymentInfo(ctx, order); err != nil {
		return fmt.Errorf("connot set order info: %w", err)
	}
	log.Info().Msg("successfully set payment order")

	if err := d.setItemInfo(ctx, order); err != nil {
		return fmt.Errorf("connot set order info: %w", err)
	}
	log.Info().Msg("successfully set item order")

	return nil
}

func (d *DataBaseOrderRepository) setOrderInfo(ctx context.Context, order entity.Order) error {
	query := `INSERT INTO "Order" (
		"order_uid", "track_number", "entry", "locale", "internal_signature", 
		"customer_id", "delivery_service", "shardkey", "sm_id", "date_created", "oof_shard"
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	return d.storage.ExecuteInsert(ctx, query, order.OrderUID, order.TrackNumber, order.Entry, order.Locale,
		order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey,
		order.SmID, order.DateCreated, order.OofShard)
}

func (d *DataBaseOrderRepository) setDeliveryInfo(ctx context.Context, order entity.Order) error {
	query := `INSERT INTO "Delivery" (
		"order_uid", "name", "phone", "zip", "city", "address", "region", "email"
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	return d.storage.ExecuteInsert(ctx, query, order.OrderUID, order.Delivery.Name, order.Delivery.Phone,
		order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
}

func (d *DataBaseOrderRepository) setPaymentInfo(ctx context.Context, order entity.Order) error {
	query := `INSERT INTO "Payment" (
		"order_uid", "transaction", "request_id", "currency", "provider", "amount", 
		"payment_dt", "bank", "delivery_cost", "goods_total", "custom_fee"
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	return d.storage.ExecuteInsert(ctx, query, order.OrderUID, order.Payment.Transaction, order.Payment.RequestID,
		order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt,
		order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
}

func (d *DataBaseOrderRepository) setItemInfo(ctx context.Context, order entity.Order) error {
	query := `INSERT INTO "Item" (
		"order_uid", "chrt_id", "track_number", "price", "rid", "name", "sale", 
		"size", "total_price", "nm_id", "brand", "status"
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	for _, item := range order.Items {
		err := d.storage.ExecuteInsert(ctx, query, order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid,
			item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return err
		}
	}

	return nil
}
