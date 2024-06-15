package domain

import (
	"context"

	"L_0_WB/internal/entity"
	"L_0_WB/internal/repository"
)

type OrderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) *OrderService {
	return &OrderService{repo}
}

func (o OrderService) Set(ctx context.Context, order entity.Order) error {
	return o.repo.Set(ctx, order)
}
