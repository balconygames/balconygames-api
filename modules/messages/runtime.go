package pixel

import (
	"github.com/go-chi/chi"
	socketio "github.com/googollee/go-socket.io"
	"github.com/kelseyhightower/envconfig"

	"gitlab.com/balconygames/analytics/modules/messages/internal/handlers"
	"gitlab.com/balconygames/analytics/modules/messages/internal/service"
	"gitlab.com/balconygames/analytics/pkg/runtime"
)

type spec struct {
	Env                 string `envconfig:"ENV" required:"True"`
	AES256Key           string `envconfig:"AES256_KEY" required:"True"`
	JWTHMACSHA256Secret string `envconfig:"JWT_SHA_256_SECRET"`
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

	svc := service.NewService(r.Logger)
	handlers.New(svc)

	// we could build in app messages using socket.io implementation
	r.WithSocketEvent("/v1/messages", "example", func(s socketio.Conn, msg string) {
	})

	r.WithRoutes(func(r1 chi.Router) {
		r.WithClientAuth(r1, func(r2 chi.Router) {
		})
	})

	return nil
}
