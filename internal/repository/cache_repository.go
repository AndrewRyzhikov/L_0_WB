package repository

import (
	"context"
	"fmt"
	"sync"

	"L_0_WB/internal/entity"
)

type CacheOrderRepository struct {
	cache map[string]entity.Order
	sync.RWMutex
	db *DataBaseOrderRepository
}

func NewCacheRepository(db *DataBaseOrderRepository, ctx context.Context) (*CacheOrderRepository, error) {
	cache := make(map[string]entity.Order)

	Orders, err := db.GetAllOrders(ctx)
	if err != nil {
		return nil, fmt.Errorf("the database could not be cached: %s", err)
	}

	for _, order := range Orders {
		cache[order.OrderUID] = order
	}

	return &CacheOrderRepository{db: db, cache: cache}, nil
}

func (repo *CacheOrderRepository) Set(ctx context.Context, order entity.Order) error {
	repo.Lock()
	repo.cache[order.OrderUID] = order
	defer repo.Unlock()

	err := repo.db.Set(ctx, order)
	if err != nil {
		return err
	}

	return nil
}

func (repo *CacheOrderRepository) Get(ctx context.Context, orderUID string) (entity.Order, error) {
	repo.RLock()
	defer repo.RUnlock()

	order, ok := repo.cache[orderUID]
	if !ok {
		return entity.Order{}, fmt.Errorf("order not found")
	}

	return order, nil
}
