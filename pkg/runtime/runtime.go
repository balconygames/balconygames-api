package runtime

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	socketio "github.com/googollee/go-socket.io"
	"github.com/pkg/errors"
	"github.com/ztrue/tracerr"
	"go.uber.org/zap"

	"gitlab.com/balconygames/analytics/pkg/auth"
	"gitlab.com/balconygames/analytics/pkg/logging"
)

var shutdownTimeout = 5 * time.Second
var requestTimeout = 60 * time.Second

type InitFunc func(*Runtime) error

type Spec struct {
	Env      string `envconfig:"ENV" required:"True"`
	HTTPPort int    `envconfig:"HTTP_PORT" default:"5000"`

	JWTClientSecret string `envconfig:"JWT_CLIENT_SHA_256_SECRET"`
	JWTServerSecret string `envconfig:"JWT_SERVER_SHA_256_SECRET"`
}

func (s Spec) Dev() bool {
	return s.Env == "dev" || s.Env == "staging" || s.Env == "test"
}

type Closable interface {
	Close()
}

type Runtime struct {
	router    chi.Router
	sockets   *socketio.Server
	spec      Spec
	closeCh   chan struct{}
	closables []Closable
	readyCh   chan struct{}

	errCh chan error
	dieCh chan os.Signal

	// action should determine what to do
	// types:
	// default - running web api
	// migrate - running migrate/migrate package
	action string

	Logger *zap.SugaredLogger
}

const apiModule = "api"

// New creates new runtime
func New(action string, s Spec) *Runtime {
	r := &Runtime{
		router:  chi.NewRouter(),
		closeCh: make(chan struct{}),
		readyCh: make(chan struct{}),
		dieCh:   make(chan os.Signal),
		errCh:   make(chan error),
		spec:    s,
		action:  action,
	}

	if r.action == dbMigrateCommand || r.action == dbResetCommand {
		// We should run only api layer, migrate is only runnable
		// once.
		return r
	}

	if r.action == "api" {
		// add handlers routing
		// A good base middleware stack
		r.router.Use(middleware.RequestID)
		r.router.Use(middleware.RealIP)
		r.router.Use(middleware.Logger)
		r.router.Use(middleware.Recoverer)

		// Set a timeout value on the request context (ctx), that will signal
		// through ctx.Done() that the request has timed out and further
		// processing should be stopped.
		r.router.Use(middleware.Timeout(requestTimeout))
	}

	return r
}

// WithSockets should starts socket server and attach to path handler
// and later it would be possible to serve it.
func (r *Runtime) WithSockets(fn func(s *socketio.Server)) {
	if r.sockets != nil {
		fn(r.sockets)
		return
	}

	// add sockets messaging layer
	server, err := socketio.NewServer(nil)
	if err != nil {
		r.errCh <- errors.Wrap(err, "socket server failed to create")
		panic(err)
	}

	r.sockets = server

	server.OnConnect("/", func(s socketio.Conn) error {
		r.Logger.
			With("id", s.ID(), "url", s.URL(), "ip", s.LocalAddr().String()).
			Info("connected user")
		return nil
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		r.Logger.
			With("id", s.ID(), "url", s.URL(), "ip", s.LocalAddr().String()).
			Infof("disconnected user reason: %s", reason)
	})

	r.WithRoutes(func(r1 chi.Router) {
		r1.Get("/v1/sockets", server.ServeHTTP)
	})

	go func() {
		err := server.Serve()
		if err != nil {
			r.errCh <- errors.Wrap(err, "socket server failed to serve")
		}
	}()
}

func (r *Runtime) WithSocketEvent(path, command string, fn func(s socketio.Conn, msg string)) {
	r.sockets.OnEvent(path, command, fn)
}

func (r *Runtime) WithSocketError(path, command string, fn func(s socketio.Conn, err error)) {
	r.sockets.OnError(path, fn)
}

func (r *Runtime) WithRoutes(fn func(r chi.Router)) {
	fn(r.router)
}

func (r *Runtime) WithClientAuth(base chi.Router, fn func(r chi.Router)) {
	base.Group(func(router chi.Router) {
		router.Use(auth.NewJWTHttpVerifierMiddleware(r.spec.JWTClientSecret))
		router.Use(auth.NewJWTUserMiddleware(r.spec.JWTClientSecret))

		fn(router)
	})
}

func (r *Runtime) WithServerAuth(base chi.Router, fn func(r chi.Router)) {
	base.Group(func(router chi.Router) {
		router.Use(auth.NewJWTHttpVerifierMiddleware(r.spec.JWTClientSecret))
		router.Use(auth.NewJWTUserMiddleware(r.spec.JWTClientSecret))

		fn(router)
	})
}

func (r *Runtime) WithClientTokenSigner(base chi.Router, fn func(r chi.Router)) {
	base.Group(func(router chi.Router) {
		router.Use(auth.NewTokenSignerMiddleware(r.spec.JWTClientSecret))

		fn(router)
	})
}

func (r *Runtime) WithServerTokenSigner(base chi.Router, fn func(r chi.Router)) {
	base.Group(func(router chi.Router) {
		router.Use(auth.NewTokenSignerMiddleware(r.spec.JWTClientSecret))

		fn(router)
	})
}

func (r *Runtime) WithLogger(env, namespace string) error {
	l, err := logging.ConfigForEnv(env).Build(
		zap.Fields(zap.String("project", namespace)),
	)
	if err != nil {
		return err
	}
	r.Logger = l.Sugar()

	return nil
}

// WithClosable registers Closable resource which will be cleaned up
// upon runtime Close call. Useful for tests to avoid db connection issues.
func (r *Runtime) WithClosable(c Closable) {
	r.closables = append(r.closables, c)
}

func health(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("OK"))
}

func (r *Runtime) Run() error {
	if r.action == "db.migrate" || r.action == "db.reset" {
		// We should run only api layer, migrate is only runnable
		// once.
		return nil
	}

	s := r.spec

	r.router.Get("/healthz", health)
	r.printRoutes()

	l, err := logging.ConfigForEnv(s.Env).Build()
	if err != nil {
		return err
	}
	defer l.Sync()
	logger := l.Sugar()

	httpListener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.HTTPPort))
	if err != nil {
		return tracerr.Wrap(errors.Wrapf(err, "http listener: could not listen to :%d.", s.HTTPPort))
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.HTTPPort),
		Handler: r.router,
	}

	signal.Notify(r.dieCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := httpServer.Serve(httpListener); err != nil && err != http.ErrServerClosed {
			r.errCh <- errors.Wrap(err, "http server failed")
		}
	}()
	logger.Infof("started http server on port: %d", s.HTTPPort)

	// broadcast that runtime is ready
	close(r.readyCh)

	// run runtime loop and wait for specific event
	// done or close
	select {
	case <-r.closeCh:
	case <-r.dieCh:
	case err := <-r.errCh:
		logger.Errorf("received error signal: %s", tracerr.Wrap(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		return tracerr.Wrap(fmt.Errorf("http server shutdown failed: %+v", err))
	}
	logger.Info("stopped http server")

	return nil
}

func (r *Runtime) Close() {
	r.closeCh <- struct{}{}
	for _, c := range r.closables {
		c.Close()
	}
}

func (r *Runtime) Wait() {
	<-r.readyCh
}

func (r *Runtime) printRoutes() {
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(r.router, walkFunc); err != nil {
		fmt.Printf("Logging err: %s\n", err.Error())
	}
}
