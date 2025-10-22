package main

import (
	"cards-service/internal/config"
	"cards-service/internal/core/app"
	"log"

	"github.com/mwinyimoha/commons/pkg/logging"
	"go.uber.org/zap"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	logger, err := logging.NewLoggerConfig().BuildLogger()
	if err != nil {
		log.Fatal("could not initialize logging", err)
	}

	defer logger.Sync()

	cfg, err := config.New()
	if err != nil {
		logger.Fatal("could not initialize app config", zap.String("original_error", err.Error()))
	}

	logger.Info("configuration loaded successfully", zap.Any("config", cfg))

	_, err = app.NewService()
	if err != nil {
		logger.Fatal("could not initialize service", zap.String("original_error", err.Error()))
	}
}
