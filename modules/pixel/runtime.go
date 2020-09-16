package pixel

import (
	"github.com/go-chi/chi"
	"github.com/kelseyhightower/envconfig"

	"gitlab.com/balconygames/analytics/modules/pixel/internal/handlers"
	"gitlab.com/balconygames/analytics/modules/pixel/internal/mq"
	"gitlab.com/balconygames/analytics/modules/pixel/internal/service"
	"gitlab.com/balconygames/analytics/pkg/nsq"
	"gitlab.com/balconygames/analytics/pkg/runtime"
)

type spec struct {
	Env                 string `envconfig:"ENV" required:"True"`
	AES256Key           string `envconfig:"AES256_KEY" required:"True"`
	JWTHMACSHA256Secret string `envconfig:"JWT_SHA_256_SECRET"`

	// NSQConfig configuration
	NSQConfig nsq.Config `envconfig:"NSQ" required:"True"`
}

// New creates pixel implementation:
// - should require to have endpoint /pixel
// - should verify authentication for the user
// - should pass the messages to message queue
func New(r *runtime.Runtime) error {
	var s spec
	if err := envconfig.Process("MODULE_PIXEL", &s); err != nil {
		return err
	}
	return withSpec(r, s)
}

func withSpec(r *runtime.Runtime, s spec) error {
	err := r.WithLogger(s.Env, "modules/auth")
	if err != nil {
		return err
	}
	defer r.Logger.Sync()

	// Pass nsq.Config
	m, err := mq.NewNSQ(r.Logger, s.NSQConfig.Topic, s.NSQConfig.Addr)
	if err != nil {
		return err
	}
	r.WithClosable(m)

	svc := service.NewService(r.Logger, m)
	h := handlers.New(svc)

	r.WithRoutes(func(r1 chi.Router) {
		r.WithClientAuth(r1, func(r2 chi.Router) {
			// Add JWT token here to verify request and associate
			// requests with their device-id, user-id
			// assume that we have user-id all the time
			// on requesting pixel.
			r2.Post("/pixel/v1/request", h.PixelHandler)

			// TODO: add GET /request and pass data as part of query params
		})
	})

	return nil
}
