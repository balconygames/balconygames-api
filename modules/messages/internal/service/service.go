package service

import (
	"go.uber.org/zap"
)

// Service contains all dependencies to perform common service tasks.
type Service struct {
	Logger *zap.SugaredLogger
}

func NewService(l *zap.SugaredLogger) *Service {
	return &Service{
		Logger: l,
	}
}
