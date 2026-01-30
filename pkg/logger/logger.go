package logger

import "go.uber.org/zap"

func InitLogger(mode string) *zap.SugaredLogger {
	var logger *zap.Logger
	var err error

	if mode == "prod" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}

	sugar := logger.Sugar()
	sugar.Infow("Logger initialized successfully", "mode", mode)

	return sugar
}
