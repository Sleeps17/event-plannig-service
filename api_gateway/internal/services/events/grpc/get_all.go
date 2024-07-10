package grpc

import (
	"context"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/domain/models"
	eventservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/events"
	eventsv1 "github.com/Sleeps17/events-planning-service-backend/api_gateway/protos/gen/go/events"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (s *Service) GetAll(ctx context.Context, startDate, andDate time.Time) ([]*models.Event, error) {
	resp, err := s.client.GetAllBetweenTwoDates(ctx, &eventsv1.GetAllBetweenTwoDatesRequest{
		StartDate: timestamppb.New(startDate),
		EndDate:   timestamppb.New(andDate),
	})

	if err != nil {
		return nil, eventservice.ErrInternalServer
	}

	var events []*models.Event
	for _, event := range resp.Events {
		var participants []*models.Employee
		for _, p := range event.Participants {
			participants = append(participants, &models.Employee{
				ID:           p.Id,
				FirstName:    p.FirstName,
				LastName:     p.LastName,
				Patronymic:   p.Patronymic,
				Email:        p.Email,
				MobileNumber: p.MobileNumber,
			})
		}

		events = append(events, &models.Event{
			ID:          event.Id,
			Title:       event.Title,
			Description: event.Description,
			StartDate:   event.StartDate.AsTime(),
			EndDate:     event.EndDate.AsTime(),
			Room: &models.Room{
				ID:       event.Room.Id,
				Name:     event.Room.Name,
				Capacity: event.Room.Capacity,
			},
			Creator: &models.Employee{
				ID:           event.Creator.Id,
				FirstName:    event.Creator.FirstName,
				LastName:     event.Creator.LastName,
				Patronymic:   event.Creator.Patronymic,
				Email:        event.Creator.Email,
				MobileNumber: event.Creator.MobileNumber,
			},
			Participants: participants,
		})
	}

	return events, nil
}
