package grpc

import (
	"context"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/domain/models"
	eventservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/events"
	eventsv1 "github.com/Sleeps17/events-planning-service-backend/api_gateway/protos/gen/go/events"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) GetByID(ctx context.Context, id uint64) (*models.Event, error) {
	resp, err := s.client.GetByID(ctx, &eventsv1.GetRequest{Id: id})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, eventservice.ErrEventNotFound
			case codes.Internal:
				return nil, eventservice.ErrEventNotFound
			}
		} else {
			return nil, eventservice.ErrInternalServer
		}
	}

	var participants []*models.Employee
	for _, p := range resp.Event.Participants {
		participants = append(participants, &models.Employee{
			ID:           p.Id,
			FirstName:    p.FirstName,
			LastName:     p.LastName,
			Patronymic:   p.Patronymic,
			Email:        p.Email,
			MobileNumber: p.MobileNumber,
		})
	}

	return &models.Event{
		ID:          resp.Event.Id,
		Title:       resp.Event.Title,
		Description: resp.Event.Description,
		StartDate:   resp.Event.StartDate.AsTime(),
		EndDate:     resp.Event.EndDate.AsTime(),
		Room: &models.Room{
			ID:       resp.Event.Room.Id,
			Name:     resp.Event.Room.Name,
			Capacity: resp.Event.Room.Capacity,
		},
		Creator: &models.Employee{
			ID:           resp.Event.Creator.Id,
			FirstName:    resp.Event.Creator.FirstName,
			LastName:     resp.Event.Creator.LastName,
			Patronymic:   resp.Event.Creator.Patronymic,
			Email:        resp.Event.Creator.Email,
			MobileNumber: resp.Event.Creator.MobileNumber,
		},
		Participants: participants,
	}, nil
}
