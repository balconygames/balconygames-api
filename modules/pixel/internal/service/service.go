package service

import (
	"go.uber.org/zap"

	"gitlab.com/balconygames/analytics/modules/pixel/internal/mq"
)

// Service contains all dependencies to perform common service tasks.
type Service struct {
	Logger *zap.SugaredLogger
	m      mq.MQable
}

func NewService(l *zap.SugaredLogger, m mq.MQable) *Service {
	return &Service{
		Logger: l,
		m:      m,
	}
}

// Push should send message to message queue
// it should be in async way handled
func (s Service) Push(msg interface{}) error {
	return s.m.Push(msg)
}
