package internalhttp

import (
	"context"
	"net"
	"net/http"

	"github.com/alexei38/otus_hw/hw12_13_14_15_calendar/internal/config"
)

type Server struct {
	Srv *http.Server
	app Application
}

type Application interface { // TODO
}

func NewServer(app Application, config config.HTTPConfig) *Server {
	s := &Server{app: app}
	s.Srv = &http.Server{
		Handler: LoggingMiddleware(http.TimeoutHandler(
			http.HandlerFunc(HelloWorld),
			config.HTTPTimeout,
			"request timeout",
		)),
		Addr:         net.JoinHostPort(config.Host, config.Port),
		WriteTimeout: config.WriteTimeout,
		ReadTimeout:  config.ReadTimeout,
		IdleTimeout:  config.IdleTimeout,
	}
	return s
}

func (s *Server) Start(ctx context.Context) error {
	return s.Srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.Srv.Shutdown(ctx)
}

// TODO
