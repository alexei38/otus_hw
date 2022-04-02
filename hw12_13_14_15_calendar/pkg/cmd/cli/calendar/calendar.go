package calendar

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexei38/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/alexei38/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/alexei38/otus_hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/alexei38/otus_hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/alexei38/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/alexei38/otus_hw/hw12_13_14_15_calendar/internal/storage/memorystorage"
	"github.com/alexei38/otus_hw/hw12_13_14_15_calendar/internal/storage/sqlstorage"
	log "github.com/sirupsen/logrus"
)

func Run() {
	config := config.NewConfig()
	err := logger.New(config.Logger)
	if err != nil {
		log.Fatalf("failed initialize logger: %v", err)
	}

	mainCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var db storage.Storage
	if config.Database.InMemory {
		db = memorystorage.New()
	} else {
		db = sqlstorage.New()
		err := db.Connect(mainCtx, config.Database.DSN)
		if err != nil {
			log.Fatalf("failed connect to database: %v", err) // nolint: gocritic
		}
	}

	calendar := app.New(db)

	server := internalhttp.NewServer(calendar, config.HTTP)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			log.Errorf("failed to stop http server: %v", err.Error())
		}
	}()

	log.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		log.Errorf("failed to start http server: %v", err.Error())
		cancel()
		os.Exit(1)
	}
}
