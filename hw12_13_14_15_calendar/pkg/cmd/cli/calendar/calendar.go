package calendar

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/alexei38/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/alexei38/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/alexei38/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/alexei38/otus_hw/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/alexei38/otus_hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/alexei38/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/alexei38/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/alexei38/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func Run() error {
	log.Infof("starting calendar")
	config, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("failed read config: %w", err)
	}
	err = logger.New(config.Logger)
	if err != nil {
		return fmt.Errorf("failed initialize logger: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var store storage.Storage
	switch config.Storage.Type {
	case "memory":
		store = memory.New()
	case "sql":
		db := sql.New()
		err := db.Connect(ctx, config.Storage.ConnectionString)
		if err != nil {
			return fmt.Errorf("failed connect to database: %v", err)

		}
		store = db
	default:
		return fmt.Errorf("unknown storage driver: %s. use sql or memory as storage type", config.Storage.Type)
	}
	calendar := app.New(store)

	errWg, errCtx := errgroup.WithContext(ctx)

	grpcServer, err := grpc.NewServer(calendar, config)
	if err != nil {
		return err
	}
	errWg.Go(func() error {
		if err := grpcServer.Start(); err != nil {
			return err
		}
		return nil
	})
	errWg.Go(func() error {
		<-errCtx.Done()
		stop()
		grpcServer.Stop()
		return nil
	})

	httpServer, err := http.NewServer(config)
	if err != nil {
		return err
	}
	errWg.Go(func() error {
		if err := httpServer.Start(errCtx); err != nil {
			return err
		}
		return nil
	})
	errWg.Go(func() error {
		<-errCtx.Done()
		stop()
		return httpServer.Stop(errCtx)
	})

	err = errWg.Wait()
	if err == context.Canceled || err == nil {
		log.Info("gracefully quit server")
		return nil
	}
	return err
}
