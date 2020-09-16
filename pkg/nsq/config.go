package nsq

type Config struct {
	// Addr host
	Addr string `envconfig:"ADDR" required:"True"`
	// Topic subscription topic
	Topic string `envconfig:"TOPIC" required:"True"`
}
