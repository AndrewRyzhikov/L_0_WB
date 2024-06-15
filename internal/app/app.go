package app

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"

	"L_0_WB/internal/config"
	"L_0_WB/internal/domain"
	"L_0_WB/internal/repository"
	"L_0_WB/internal/repository/postgres"
	"L_0_WB/internal/transport"
)

func StartApp() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config file")
	}

	log.Logger, err = setupLog(&cfg.LogConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to setup logger")
	}

	storage, err := postgres.NewStorage(cfg.PostgresConfig)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	log.Info().Msg("Database connected...")

	dataBaseRepo := repository.NewDataBaseOrderRepository(storage)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.CachingTTl)
	defer cancel()

	cacheRepo, err := repository.NewCacheRepository(dataBaseRepo, ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create cache repository")
	}

	service := domain.NewOrderService(cacheRepo)

	listener := transport.NewOrderListener(service, cfg.NatsConfig)

	err = listener.Listen()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
	}

	log.Info().Msg("Listener started...")
	log.Info().Msg("App started...")

	_, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	done := make(chan struct{})
	go func() {
		defer close(done)

		if err = listener.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close listener")
		}
		log.Info().Msg("Listener stopped...")

		if err = storage.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close database")
		}
		log.Info().Msg("Database disconnected...")
	}()

	select {
	case <-done:
		log.Info().Msg("App stopped...")
	case <-ctx.Done():
		log.Fatal().Err(ctx.Err()).Msg("Shutdown timed out")
	}
}

func setupLog(config *config.LogConfig) (zerolog.Logger, error) {
	lumberjackCfg := config.Lumberjack

	lr := &lumberjack.Logger{
		Filename:   config.Path,
		MaxSize:    int(lumberjackCfg.MaxSize),
		MaxAge:     int(lumberjackCfg.MaxAge),
		MaxBackups: int(lumberjackCfg.MaxBackups),
		LocalTime:  lumberjackCfg.LocalTime,
		Compress:   lumberjackCfg.Compress,
	}

	level, err := zerolog.ParseLevel(config.Level)

	if err != nil {
		return zerolog.Logger{}, fmt.Errorf("invalid log level: %s", config.Level)
	}

	return zerolog.New(lr).Level(level).With().Timestamp().Logger(), nil
}
