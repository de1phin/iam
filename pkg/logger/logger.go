package logger

import (
	"log"

	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {

	debugMode := "dev"

	var localLogger *zap.Logger
	var err error

	switch debugMode {
	case "prod":
		cfg := zap.NewProductionConfig()
		cfg.DisableCaller = true
		cfg.DisableStacktrace = true
		cfg.OutputPaths = []string{"stdout", "../../pkg/logs/data/log.txt"}
		localLogger, err = cfg.Build()
	default:
		localLogger, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Fatal("logger init: ", err)
	}

	logger = localLogger

	return
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}
