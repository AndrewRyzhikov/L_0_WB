package repository

import (
	"context"

	"L_0_WB/internal/entity"
)

type OrderRepository interface {
	Set(context.Context, entity.Order) error
	Get(context.Context, string) (entity.Order, error)
}
