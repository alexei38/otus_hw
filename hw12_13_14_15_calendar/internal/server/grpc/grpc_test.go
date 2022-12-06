package grpc

import (
	"context"
	"net"
	"os"
	"time"

	"github.com/alexei38/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/alexei38/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/alexei38/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/alexei38/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SuiteTest struct {
	suite.Suite
	client   EventsClient
	conn     *grpc.ClientConn
	grpcSrv  *grpc.Server
	listener *bufconn.Listener
	app      *app.App
	db       storage.Storage
}

func (s *SuiteTest) SetupTest() {
	ctx := context.Background()
	dbConnect := os.Getenv("PQ_TEST")
	var stor storage.Storage
	if dbConnect == "" {
		stor = memory.New()
	} else {
		db := sql.New()
		err := db.Connect(ctx, dbConnect)
		require.NoError(s.T(), err)
		stor = db
	}
	s.app = app.New(stor)

	s.conn, _ = grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer(s)))
	s.client = NewEventsClient(s.conn)

	// _ = s.app.DeleteAll(ctx)
}

func (s *SuiteTest) NewCommonEvent() *Event {
	var eventStart = time.Now().Add(2 * time.Hour)
	var eventStop = eventStart.Add(time.Hour)
	notification := 4 * time.Hour

	return &Event{
		Id:          0,
		Title:       "some event",
		Start:       timestamppb.New(eventStart),
		Stop:        timestamppb.New(eventStop),
		Description: "the event",
		UserID:      1,
		BeforeSend:  timestamppb.New(time.Now().Add(notification)),
	}
}

func dialer(s *SuiteTest) func(context.Context, string) (net.Conn, error) {
	s.listener = bufconn.Listen(1024 * 1024)

	s.grpcSrv = grpc.NewServer()
	RegisterEventsServer(s.grpcSrv, NewService(s.app))

	go func() {
		_ = s.grpcSrv.Serve(s.listener)
	}()

	return func(context.Context, string) (net.Conn, error) {
		return s.listener.Dial()
	}
}

func (s *SuiteTest) EqualEvents(event1, event2 *Event) {
	s.Require().Equal(event1.Title, event2.Title)
	s.Require().Equal(event1.Description, event2.Description)
	s.Require().Equal(event1.Start.AsTime().Unix(), event2.Start.AsTime().Unix())
	s.Require().Equal(event1.Stop.AsTime().Unix(), event2.Stop.AsTime().Unix())
	s.Require().Equal(event1.UserID, event2.UserID)
	if event1.BeforeSend == nil || event2.BeforeSend == nil {
		s.Require().Equal(event1.BeforeSend, event2.BeforeSend)
	} else {
		s.Require().Equal(event1.BeforeSend.AsTime(), event2.BeforeSend.AsTime())
	}
}
