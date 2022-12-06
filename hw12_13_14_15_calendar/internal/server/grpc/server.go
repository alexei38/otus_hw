package grpc

import (
	"net"

	"github.com/alexei38/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/alexei38/otus_hw/hw12_13_14_15_calendar/internal/config"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type server struct {
	app      *app.App
	srv      *grpc.Server
	listener net.Listener
}

func NewServer(app *app.App, config *config.Config) (*server, error) {
	grpcHostPort := net.JoinHostPort(config.GRPC.Host, config.GRPC.Port)
	grpcListen, err := net.Listen("tcp", grpcHostPort)
	if err != nil {
		return nil, err
	}
	return &server{
		app:      app,
		srv:      grpc.NewServer(),
		listener: grpcListen,
	}, nil
}

func (s *server) Start() error {
	log.Infof("grpc server started %v", s.listener.Addr())
	RegisterEventsServer(s.srv, NewService(s.app))
	return s.srv.Serve(s.listener)
}

func (s *server) Stop() {
	s.srv.GracefulStop()
}
