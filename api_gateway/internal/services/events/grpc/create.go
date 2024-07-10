package grpc

import (
	"context"
	"encoding/json"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/domain/models"
	eventservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/events"
	eventsv1 "github.com/Sleeps17/events-planning-service-backend/api_gateway/protos/gen/go/events"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Service) Create(ctx context.Context, event *models.Event) (*eventservice.CreateResponse, error) {
	var participants []*eventsv1.Employee
	for _, p := range event.Participants {
		participants = append(participants, &eventsv1.Employee{
			Id: p.ID,
		})
	}

	metaCtx := metadata.NewOutgoingContext(ctx, metadata.New(nil))

	resp, err := s.client.Create(metaCtx, &eventsv1.CreateRequest{
		Event: &eventsv1.Event{
			Id:          event.ID,
			Title:       event.Title,
			Description: event.Description,
			StartDate:   timestamppb.New(event.StartDate),
			EndDate:     timestamppb.New(event.EndDate),
			Room: &eventsv1.Room{
				Id: event.Room.ID,
			},
			Creator: &eventsv1.Employee{
				Id: event.Creator.ID,
			},
			Participants: participants,
		},
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.AlreadyExists:
				return nil, eventservice.ErrRoomIsOccupied
			case codes.FailedPrecondition:
				header, exist := metadata.FromIncomingContext(metaCtx)
				if !exist {
					return nil, eventservice.ErrInternalServer
				}

				jsonBusyEmployees := header["busy-employees"]
				var busyEmployees []*models.Employee
				if len(jsonBusyEmployees) > 0 {
					if err := json.Unmarshal(
						[]byte(jsonBusyEmployees[0]),
						&busyEmployees,
					); err != nil {
						return nil, eventservice.ErrInternalServer
					}
				} else {
					return nil, eventservice.ErrInternalServer
				}

				return &eventservice.CreateResponse{
					BusyEmployees: busyEmployees,
				}, eventservice.ErrSomeWorkersAreBusy
			case codes.Internal:
				return nil, eventservice.ErrInternalServer
			}
		} else {
			return nil, eventservice.ErrInternalServer
		}
	}

	return &eventservice.CreateResponse{
		ID: resp.Id,
	}, nil
}
