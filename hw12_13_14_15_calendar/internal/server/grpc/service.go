package grpc

import (
	"context"

	"github.com/alexei38/otus_hw/hw12_13_14_15_calendar/internal/app"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	UnimplementedEventsServer
	app *app.App
}

func NewService(app *app.App) *Service {
	return &Service{
		app: app,
	}
}

func (s *Service) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	id, err := s.app.Create(
		ctx,
		int(req.Event.Id),
		req.Event.Title,
		req.Event.Start.AsTime(),
		req.Event.Stop.AsTime(),
		req.Event.Description,
		int(req.Event.UserID),
		req.Event.BeforeSend.AsTime(),
	)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &CreateResponse{EventId: int32(id)}, nil
}
