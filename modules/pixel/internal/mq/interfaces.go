package mq

// MQable should be implemented for different implementations
// of message queue pusher.
type MQable interface {
	Push(message interface{}) error
}
