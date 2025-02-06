package internalgrpc

import (
	"context"

	"github.com/evg555/hw-otus/hw12_13_14_15_calendar/api/pb"
	"github.com/evg555/hw-otus/hw12_13_14_15_calendar/internal/app"
)

type Handler struct {
	app    Application
	logger Logger

	pb.UnsafeEventServiceServer
}

func (h Handler) CreateEvent(ctx context.Context, req *pb.CreateRequest) (*pb.Response, error) {
	err := h.app.CreateEvent(ctx, req.GetId(), req.GetTitle())
	if err != nil {
		h.logger.Error(err.Error())
		return renderErrorResponse(err), err
	}

	return &pb.Response{}, nil
}

func (h Handler) UpdateEvent(ctx context.Context, req *pb.UpdateRequest) (*pb.Response, error) {
	id := req.GetId()

	event := app.Event{
		ID:          id,
		Title:       req.GetEvent().GetTitle(),
		StartDate:   req.GetEvent().GetStartDate(),
		EndDate:     req.GetEvent().GetEndDate(),
		Description: req.GetEvent().GetDescription(),
		UserID:      req.GetEvent().GetUserId(),
		NotifyDays:  req.GetEvent().GetNotifyDays(),
	}

	err := h.app.UpdateEvent(ctx, id, event)
	if err != nil {
		h.logger.Error(err.Error())
		return renderErrorResponse(err), err
	}

	return &pb.Response{}, nil
}

func (h Handler) DeleteEvent(ctx context.Context, req *pb.DeleteRequest) (*pb.Response, error) {
	id := req.GetId()

	err := h.app.DeleteEvent(ctx, id)
	if err != nil {
		h.logger.Error(err.Error())
		return renderErrorResponse(err), err
	}

	return &pb.Response{}, nil
}

func (h Handler) ListEvents(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	events, err := h.app.ListEvents(ctx, req.GetDate(), req.GetPeriod())
	if err != nil {
		h.logger.Error(err.Error())
		return &pb.ListResponse{Resp: renderErrorResponse(err)}, err
	}

	resp := pb.ListResponse{
		Resp:   &pb.Response{},
		Events: make([]*pb.Event, 0, len(events)),
	}

	for _, event := range events {
		resp.Events = append(resp.Events, &pb.Event{
			Id:          event.ID,
			Title:       event.Title,
			StartDate:   event.StartDate,
			EndDate:     event.EndDate,
			Description: event.Description,
			UserId:      event.UserID,
			NotifyDays:  event.NotifyDays,
		})
	}

	return &resp, nil
}

func renderErrorResponse(err error) *pb.Response {
	return &pb.Response{
		Error:   true,
		Message: err.Error(),
	}
}
