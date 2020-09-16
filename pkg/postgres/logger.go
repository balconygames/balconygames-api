package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

type Logger struct {
	logger *zap.SugaredLogger
}

func NewLogger(logger *zap.SugaredLogger) *Logger {
	return &Logger{logger: logger}
}

func (pl *Logger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	scoped := pl.logger

	for k, v := range data {
		scoped = scoped.With(k, v)
	}

	switch level {
	case pgx.LogLevelTrace:
		scoped = scoped.With("PGX_LOG_LEVEL", level)
		scoped.Debug(msg)
	case pgx.LogLevelDebug:
		scoped.Debug(msg)
	case pgx.LogLevelInfo:
		scoped.Info(msg)
	case pgx.LogLevelWarn:
		scoped.Warn(msg)
	case pgx.LogLevelError:
		scoped.Error(msg)
	default:
		scoped = scoped.With("PGX_LOG_LEVEL", level)
		scoped.Error(msg)
	}
}
