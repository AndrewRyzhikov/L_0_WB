package transport

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog/log"

	"L_0_WB/internal/config"
	"L_0_WB/internal/domain"
	"L_0_WB/internal/entity"
)

type OrderListener struct {
	service      *domain.OrderService
	subscription stan.Subscription
	config       config.NatsConfig
}

func NewOrderListener(service *domain.OrderService, config config.NatsConfig) *OrderListener {
	return &OrderListener{service: service, config: config}
}

func (l *OrderListener) setOrder() func(msg *stan.Msg) {
	return func(msg *stan.Msg) {
		order := entity.Order{}
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Error().Err(err).Msg("Error unmarshalling order")
		}

		ctx, cancel := context.WithTimeout(context.Background(), l.config.TTl)
		defer cancel()

		err = l.service.Set(ctx, order)
		if err != nil {
			log.Error().Err(err).Msg("Error setting order")
		}
	}
}

func (l *OrderListener) Listen() error {

	fmt.Println(l.config.StanClusterId, l.config.ClientId, l.config.ChannelName)

	sc, err := stan.Connect(l.config.StanClusterId, l.config.ClientId)
	if err != nil {
		return fmt.Errorf("stan connect error: %s", err)
	}

	s, err := sc.Subscribe(l.config.ChannelName, l.setOrder())
	if err != nil {
		return fmt.Errorf("subscribe error: %s", err)
	}

	l.subscription = s

	return nil
}

func (l *OrderListener) Close() error {
	return l.subscription.Close()
}
