package mq

import (
	"encoding/json"

	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"
)

// NSQ - implementation of nsq client to push messages
// on pixel request.
type NSQ struct {
	topic    string
	producer *nsq.Producer

	logger *zap.SugaredLogger
}

type nsqLogger struct {
	logger *zap.SugaredLogger
}

// Output requires to be impelemnted to have it as part
// of NSQ client.
func (m nsqLogger) Output(calldepth int, s string) error {
	m.logger.Debug(s)
	return nil
}

// NewNSQ should message queue implementation.
func NewNSQ(l *zap.SugaredLogger, t string, addr string) (*NSQ, error) {
	config := nsq.NewConfig()

	lScoped := l.With("scope", "mq")

	w, err := nsq.NewProducer(addr, config)
	if err != nil {
		return nil, err
	}
	w.SetLogger(&nsqLogger{lScoped}, nsq.LogLevelInfo)

	instance := &NSQ{
		producer: w,
		topic:    t,
		logger:   lScoped,
	}

	return instance, nil
}

// Push message to NSQ
func (m *NSQ) Push(message interface{}) error {
	m.logger.Debugf("marshal message to json: %v", message)

	// message could be the different struct
	// but we are passing json all the time.
	b, err := json.Marshal(message)
	if err != nil {
		return err
	}

	m.logger.Debugf("pushing the new message %v", message)

	return m.producer.PublishAsync(m.topic, b, nil)
}

// Close message queue client
func (m *NSQ) Close() {
	m.Close()
}
