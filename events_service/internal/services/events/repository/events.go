package repository_events_service

import (
	"context"
	"errors"
	"github.com/Sleeps17/events-planning-service-backend/events_service/internal/domain/models"
	eventsrepository "github.com/Sleeps17/events-planning-service-backend/events_service/internal/repository/events"
	eventservice "github.com/Sleeps17/events-planning-service-backend/events_service/internal/services/events"
	"time"
)

type Service struct {
	repo eventsrepository.Repository
}

func New(repo eventsrepository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetByID(ctx context.Context, id uint64) (*models.Event, error) {
	// Получение мероприятия
	event, err := s.repo.SelectByID(ctx, id)
	if err != nil {
		// Если мероприятие не найдено, то возращаем специальную ошибку
		if errors.Is(err, eventsrepository.ErrEventNotFound) {
			return nil, eventservice.ErrEventNotFound
		}

		// Иначе, просто ошибку
		return nil, err
	}

	return event, nil
}

func (s *Service) GetAllOfTheWeek(ctx context.Context, startDate, endDate time.Time) ([]*models.Event, error) {
	// Получение мероприятий
	events, err := s.repo.SelectAllBetweenTwoDates(ctx, startDate, endDate)
	// Если произошла ошибка, то возвращаем ее
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Service) Add(ctx context.Context, event *models.Event) (*eventservice.AddResponse, error) {
	// Проверка доступности комнаты
	roomIsAvailable, err := s.repo.CheckRoomIsAvailable(
		ctx,
		event.Room.ID,
		event.StartDate,
		event.EndDate,
	)
	if err != nil {
		return nil, err
	}
	if !roomIsAvailable {
		return nil, eventservice.ErrRoomIsNotAvailable
	}

	// Проверка доступности всех рабоотников и создателя
	employeesIDs := make([]uint64, len(event.Participants))
	for i, participant := range event.Participants {
		employeesIDs[i] = participant.ID
	}
	employeesIDs = append(employeesIDs, event.Creator.ID)
	busyEmployees, err := s.repo.CheckEmployeesAreAvailable(
		ctx,
		employeesIDs,
		event.StartDate,
		event.EndDate,
	)
	if err != nil {
		return nil, err
	}
	if len(busyEmployees) > 0 {
		return &eventservice.AddResponse{
			ID:            0,
			BusyEmployees: busyEmployees,
		}, eventservice.ErrEmployeesAreNotAvailable
	}

	// Добавление мероприятия
	id, err := s.repo.Insert(ctx, event)
	if err != nil {
		return nil, err
	}

	return &eventservice.AddResponse{
		ID:            id,
		BusyEmployees: nil,
	}, nil
}

func (s *Service) UpdateByID(ctx context.Context, updatedEvent *models.Event) (*eventservice.UpdateResponse, error) {
	// Проверка доступности комнаты
	roomIsAvailable, err := s.repo.CheckRoomIsAvailable(
		ctx,
		updatedEvent.Room.ID,
		updatedEvent.StartDate,
		updatedEvent.EndDate,
	)
	if err != nil {
		return nil, err
	}
	if !roomIsAvailable {
		return nil, eventservice.ErrRoomIsNotAvailable
	}

	// Проверка доступности всех сотрудников и создателя
	employeesIDs := make([]uint64, len(updatedEvent.Participants))
	for i, participant := range updatedEvent.Participants {
		employeesIDs[i] = participant.ID
	}
	employeesIDs = append(employeesIDs, updatedEvent.Creator.ID)
	busyEmployees, err := s.repo.CheckEmployeesAreAvailable(
		ctx,
		employeesIDs,
		updatedEvent.StartDate,
		updatedEvent.EndDate,
	)
	if err != nil {
		return nil, err
	}
	if len(busyEmployees) > 0 {
		return &eventservice.UpdateResponse{
			UpdatedEvent:  nil,
			BusyEmployees: busyEmployees,
		}, eventservice.ErrEmployeesAreNotAvailable
	}

	// Обновления данных о мероприятии
	if err := s.repo.UpdateByID(ctx, updatedEvent); err != nil {
		if errors.Is(err, eventsrepository.ErrEventNotFound) {
			return nil, eventservice.ErrEventNotFound
		}
		return nil, err
	}

	// Получение обновленного мероприятия
	event, err := s.repo.SelectByID(ctx, updatedEvent.ID)
	if err != nil {
		return nil, err
	}

	return &eventservice.UpdateResponse{UpdatedEvent: event, BusyEmployees: nil}, nil
}

func (s *Service) DeleteByID(ctx context.Context, id uint64) error {
	// Удаление мероприятия
	if err := s.repo.DeleteByID(ctx, id); err != nil {
		// Если нет мероприятия с таким id, то возвращаем специальную ошибку
		if errors.Is(err, eventsrepository.ErrEventNotFound) {
			return eventservice.ErrEventNotFound
		}
		// Иначе просто возврщаем ошибку
		return err
	}

	return nil
}
