package transport

import (
	"L_0_WB/internal/config"
	"L_0_WB/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/http"
)

type OrderController struct {
	server  *http.Server
	service *domain.OrderService
	config  config.HttpServerConfig
}

func NewOrderController(service *domain.OrderService, config config.HttpServerConfig) *OrderController {
	return &OrderController{service: service, config: config}
}

func (controller *OrderController) orderHandler(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("order handler called")

	vars := mux.Vars(r)
	uid := vars["UID"]

	if uid == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), controller.config.TTL)
	defer cancel()

	order, err := controller.service.Get(ctx, uid)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	data, err := json.Marshal(order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func (controller *OrderController) Start() error {
	m := &mux.Router{}

	m.
		Path("/order/{UID}").
		Methods(http.MethodGet).
		HandlerFunc(controller.orderHandler)

	controller.server = &http.Server{
		Handler: m,
		Addr:    controller.config.Port,
	}

	go func() {
		err := controller.server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Failed to started http Controller")
		}
	}()

	return nil
}

func (controller *OrderController) Shutdown() error {
	if err := controller.server.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("failed to stopped Healtcontrollerhecker : %w", err)
	}
	return nil
}
