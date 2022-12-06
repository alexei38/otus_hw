package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/alexei38/otus_hw/hw12_13_14_15_calendar/internal/config"
	pb "github.com/alexei38/otus_hw/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	srv         *http.Server
	grpcHost    string
	grpcPort    string
	grpcHandler *runtime.ServeMux
}

func NewServer(config *config.Config) (*Server, error) {
	httpHostPort := net.JoinHostPort(config.HTTP.Host, config.HTTP.Port)
	grpcHandler := runtime.NewServeMux()
	s := &Server{
		grpcHandler: grpcHandler,
		grpcHost:    config.GRPC.Host,
		grpcPort:    config.GRPC.Port,
		srv: &http.Server{
			Handler: LoggingMiddleware(http.TimeoutHandler(
				grpcHandler,
				config.HTTP.HTTPTimeout,
				"request timeout",
			)),
			Addr:         httpHostPort,
			WriteTimeout: config.HTTP.WriteTimeout,
			ReadTimeout:  config.HTTP.ReadTimeout,
			IdleTimeout:  config.HTTP.IdleTimeout,
		},
	}
	return s, nil
}

func (s *Server) Start(ctx context.Context) error {
	grpcHostPort := net.JoinHostPort(s.grpcHost, s.grpcPort)
	grpcClientCtx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	conn, err := grpc.DialContext(
		grpcClientCtx,
		grpcHostPort,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed connect to grpc server: %w", err)
	}
	if err = pb.RegisterEventsHandler(ctx, s.grpcHandler, conn); err != nil {
		return fmt.Errorf("failed to register gateway: %w", err)
	}
	log.Infof("http server started %v", s.srv.Addr)
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

//
// func NewServer(config *config.Config) {
// 	grpcHostPort := net.JoinHostPort(config.GRPC.Host, config.GRPC.Port)
// 	httpHostPort := net.JoinHostPort(config.HTTP.Host, config.HTTP.Port)
//
// 	grpcClientCtx, cancel := context.WithTimeout(context.Background(), time.Second*2)
// 	defer cancel()
// 	conn, err := grpc.DialContext(
// 		grpcClientCtx,
// 		grpcHostPort,
// 		grpc.WithBlock(),
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 	)
// 	if err != nil {
// 		log.Fatalln("Failed to dial server:", err)
// 	}
// 	grpcHandler := runtime.NewServeMux()
// 	err = pb.RegisterEventsHandler(context.Background(), grpcHandler, conn)
// 	if err != nil {
// 		log.Fatalln("Failed to register gateway:", err)
// 	}
//
// 	// s := &Server{app: app}
// 	// s.Srv = &http.Server{
// 	// 	Handler: LoggingMiddleware(http.TimeoutHandler(
// 	// 		grpcHandler,
// 	// 		config.HTTP.HTTPTimeout,
// 	// 		"request timeout",
// 	// 	)),
// 	// 	Addr: httpHostPort,
// 	// 	WriteTimeout: config.WriteTimeout,
// 	// 	ReadTimeout:  config.ReadTimeout,
// 	// 	IdleTimeout:  config.IdleTimeout,
// 	// }
// 	// // grpcSrv := grpc.NewServer()
// 	// var muxServer grpc.ServiceRegistrar
// 	// muxServer = grpc.NewServer()
// 	// pb.RegisterEventsServer(muxServer, pb.Service{})
// 	// mux rpc
// 	// muxServer.RegisterCodec(json.NewCodec(), "application/json")
// 	// muxServer.RegisterService(new(api.Event), "")
// 	// ("/grpc", muxServer)
// }
//
// // func (s *Server) Start(ctx context.Context) error {
// // 	return s.Srv.ListenAndServe()
// // }
// //
// // func (s *Server) Stop(ctx context.Context) error {
// // 	return s.Srv.Shutdown(ctx)
// // }

// TODO
