package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ConfigForEnv(env string) zap.Config {
	var config zap.Config
	switch env {
	case "dev":
		config = zap.NewDevelopmentConfig()
	case "staging":
		config = zap.NewDevelopmentConfig()
	case "test":
		config = zap.NewDevelopmentConfig()
	default:
		config = zap.NewProductionConfig()
	}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return config
}
