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

func (s *Service) Update(ctx context.Context, event *models.Event) (*eventservice.UpdateResponse, error) {
	var participants []*eventsv1.Employee
	for _, participant := range event.Participants {
		participants = append(participants, &eventsv1.Employee{
			Id: participant.ID,
		})
	}

	metaCtx := metadata.NewOutgoingContext(ctx, metadata.New(nil))

	resp, err := s.client.Update(metaCtx, &eventsv1.UpdateRequest{
		Id: event.ID,
		UpdatedEvent: &eventsv1.Event{
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
			case codes.NotFound:
				return nil, eventservice.ErrEventNotFound
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

				return &eventservice.UpdateResponse{
					BusyEmployees: busyEmployees,
				}, eventservice.ErrSomeWorkersAreBusy
			case codes.Internal:
				return nil, eventservice.ErrInternalServer
			}
		} else {
			return nil, eventservice.ErrInternalServer
		}
	}

	var p []*models.Employee
	for _, participant := range resp.UpdatedEvent.Participants {
		p = append(p, &models.Employee{
			ID:           participant.Id,
			FirstName:    participant.FirstName,
			LastName:     participant.LastName,
			Patronymic:   participant.Patronymic,
			Email:        participant.Email,
			MobileNumber: participant.MobileNumber,
		})
	}

	return &eventservice.UpdateResponse{
		UpdatedEvent: &models.Event{
			ID:          resp.UpdatedEvent.Id,
			Title:       resp.UpdatedEvent.Title,
			Description: resp.UpdatedEvent.Description,
			StartDate:   resp.UpdatedEvent.StartDate.AsTime(),
			EndDate:     resp.UpdatedEvent.EndDate.AsTime(),
			Room: &models.Room{
				ID:       resp.UpdatedEvent.Room.Id,
				Name:     resp.UpdatedEvent.Room.Name,
				Capacity: resp.UpdatedEvent.Room.Capacity,
			},
			Creator: &models.Employee{
				ID:           resp.UpdatedEvent.Creator.Id,
				FirstName:    resp.UpdatedEvent.Creator.FirstName,
				LastName:     resp.UpdatedEvent.Creator.LastName,
				Patronymic:   resp.UpdatedEvent.Creator.Patronymic,
				Email:        resp.UpdatedEvent.Creator.Email,
				MobileNumber: resp.UpdatedEvent.Creator.MobileNumber,
			},
			Participants: p,
		},
	}, nil
}
