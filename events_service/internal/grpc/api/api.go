package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Sleeps17/event-plannig-service-backend/events-service/internal/domain/models"
	employeeservice "github.com/Sleeps17/event-plannig-service-backend/events-service/internal/services/employees"
	eventservice "github.com/Sleeps17/event-plannig-service-backend/events-service/internal/services/events"
	roomservice "github.com/Sleeps17/event-plannig-service-backend/events-service/internal/services/rooms"
	eventsv1 "github.com/Sleeps17/event-plannig-service-backend/events-service/protos/gen/go/events"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
)

type api struct {
	eventsv1.UnimplementedEventsServer
	logger            *slog.Logger
	eventProvider     eventservice.Service
	roomsProvider     roomservice.Service
	employeesProvider employeeservice.Service
}

func Register(
	srv *grpc.Server,
	logger *slog.Logger,
	eventProvider eventservice.Service,
	roomsProvider roomservice.Service,
	employeesProvider employeeservice.Service,
) {
	eventsv1.RegisterEventsServer(srv, &api{
		logger:            logger,
		eventProvider:     eventProvider,
		roomsProvider:     roomsProvider,
		employeesProvider: employeesProvider,
	})
}

func (a *api) Create(ctx context.Context, request *eventsv1.CreateRequest) (*eventsv1.CreateResponse, error) {
	event := transformEventFromRequest(request.GetEvent())

	a.logger.Info("try to handle create request", slog.Any("event", event))

	resp, err := a.eventProvider.Add(ctx, event)
	if err != nil {
		// Если комната занята
		if errors.Is(err, eventservice.ErrRoomIsNotAvailable) {
			a.logger.Info("room is not available", slog.Any("roomID", event.Room.ID))
			return nil, status.Error(codes.AlreadyExists, RoomIsNotAvailableMsg)
		}

		// Если какой-то из работников занят
		if errors.Is(err, eventservice.ErrEmployeesAreNotAvailable) {

			a.logger.Info("employees are not available", slog.Any("busy employees", resp.BusyEmployees))

			// Обогащаем данные о занятых сотрудниках
			aggregatedBusyEmployees, err := a.getInfoAboutBusyEmployees(ctx, resp.BusyEmployees)
			if err != nil {
				a.logger.Error("failed to get info about busy employees", slog.Any("error", err))
				return nil, status.Error(codes.Internal, InternalErrorMsg)
			}

			// Преобразуем их в строку
			busyEmployeesStr, err := json.Marshal(aggregatedBusyEmployees)
			if err != nil {
				a.logger.Error("failed to marshal busy employees", slog.Any("error", err))
				return nil, status.Error(codes.Internal, InternalErrorMsg)
			}

			// Кладем в метадату
			md := metadata.Pairs("busy-employees", string(busyEmployeesStr))
			if err := grpc.SetHeader(ctx, md); err != nil {
				a.logger.Error("failed to set header", slog.Any("error", err))
				return nil, status.Error(codes.Internal, InternalErrorMsg)
			}

			// Возвращаем ошибку
			return nil, status.Error(codes.FailedPrecondition, EmployeesAreNotAvailableMsg)
		}

		// Если произошла ошибка
		a.logger.Error("failed to add event", slog.Any("error", err))
		return nil, status.Error(codes.Internal, InternalErrorMsg)
	}

	a.logger.Info("successfully created event", slog.Any("eventID", resp.ID))

	return &eventsv1.CreateResponse{Id: resp.ID}, nil
}

func (a *api) GetByID(ctx context.Context, request *eventsv1.GetRequest) (*eventsv1.GetResponse, error) {
	id := request.GetId()

	a.logger.Info("try to get event by id", slog.Any("id", id))

	// Получаем мероприятие
	event, err := a.eventProvider.GetByID(ctx, id)
	if err != nil {
		a.logger.Error("failed to get event by id", slog.Any("id", id), slog.Any("error", err))
		return nil, status.Error(codes.Internal, InternalErrorMsg)
	}

	// Обогащаем данные о аудитории
	event.Room, err = a.getInfoAboutRoom(ctx, event.Room.ID)
	if err != nil {
		a.logger.Error("failed to get room info", slog.Any("id", event.Room.ID), slog.Any("error", err))
		return nil, status.Error(codes.Internal, InternalErrorMsg)
	}

	// Обогащаем данные об участниках
	event.Participants, err = a.getInfoAboutParticipants(ctx, event.Participants)
	if err != nil {
		a.logger.Error("failed to get participants info", slog.Any("id", event.ID), slog.Any("error", err))
		return nil, status.Error(codes.Internal, InternalErrorMsg)
	}

	// Обогащаем данные о создателе
	event.Creator, err = a.getInfoAboutCreator(ctx, event.Creator.ID)
	if err != nil {
		a.logger.Error("failed to get creator info", slog.Any("id", event.Creator.ID), slog.Any("error", err))
		return nil, status.Error(codes.Internal, InternalErrorMsg)
	}

	a.logger.Info("successfully get event by id", slog.Any("event", event))
	return &eventsv1.GetResponse{Event: transformEventToResponse(event)}, nil
}

func (a *api) GetAllBetweenTwoDates(ctx context.Context, request *eventsv1.GetAllBetweenTwoDatesRequest) (*eventsv1.GetAllBetweenTwoDatesResponse, error) {
	startDate := request.GetStartDate().AsTime()
	endDate := request.GetEndDate().AsTime()

	a.logger.Info("try to get all events of the week", slog.Any("start_date", startDate), slog.Any("end_date", endDate))

	events, err := a.eventProvider.GetAllOfTheWeek(ctx, startDate, endDate)
	if err != nil {
		a.logger.Error("failed to get all events of the week", slog.Any("error", err))
		return nil, status.Error(codes.Internal, InternalErrorMsg)
	}

	for i := range events {
		// Обогащаем данные о аудитории
		events[i].Room, err = a.getInfoAboutRoom(ctx, events[i].Room.ID)
		if err != nil {
			a.logger.Error("failed to get room info", slog.Any("id", events[i].Room.ID), slog.Any("error", err))
			return nil, status.Error(codes.Internal, InternalErrorMsg)
		}

		// Обогащаем данные об участниках
		events[i].Participants, err = a.getInfoAboutParticipants(ctx, events[i].Participants)
		if err != nil {
			a.logger.Error("failed to get participants info", slog.Any("id", events[i].ID), slog.Any("error", err))
			return nil, status.Error(codes.Internal, InternalErrorMsg)
		}

		// Обогащаем данные о создателе
		events[i].Creator, err = a.getInfoAboutCreator(ctx, events[i].Creator.ID)
		if err != nil {
			a.logger.Error("failed to get creator info", slog.Any("id", events[i].Creator.ID), slog.Any("error", err))
			return nil, status.Error(codes.Internal, InternalErrorMsg)
		}
	}

	eventsForResponse := make([]*eventsv1.Event, len(events))
	for i, event := range events {
		eventsForResponse[i] = transformEventToResponse(event)
	}

	a.logger.Info("successfully get all events of the week", slog.Any("start_date", startDate), slog.Any("end_date", endDate))

	return &eventsv1.GetAllBetweenTwoDatesResponse{Events: eventsForResponse}, nil
}

func (a *api) Update(ctx context.Context, request *eventsv1.UpdateRequest) (*eventsv1.UpdateResponse, error) {
	event := transformEventFromRequest(request.GetUpdatedEvent())

	a.logger.Info("try to update event", slog.Any("eventID", event.ID))

	resp, err := a.eventProvider.UpdateByID(ctx, event)
	if err != nil {
		// Если комната занята
		if errors.Is(err, eventservice.ErrRoomIsNotAvailable) {
			a.logger.Info("room is not available", slog.Any("roomID", event.Room.ID))
			return nil, status.Error(codes.AlreadyExists, RoomIsNotAvailableMsg)
		}

		// Если какой-то из работников занят
		if errors.Is(err, eventservice.ErrEmployeesAreNotAvailable) {
			a.logger.Info("employees are not available", slog.Any("busy employees", resp.BusyEmployees))

			// Обогащаем данные о занятых сотрудниках
			aggregatedBusyEmployees, err := a.getInfoAboutBusyEmployees(ctx, resp.BusyEmployees)
			if err != nil {
				a.logger.Error("failed to get info about busy employees", slog.Any("error", err))
				return nil, status.Error(codes.Internal, InternalErrorMsg)
			}

			// Преобразуем их в строку
			busyEmployeesStr, err := json.Marshal(aggregatedBusyEmployees)
			if err != nil {
				a.logger.Error("failed to marshal busy employees", slog.Any("error", err))
				return nil, status.Error(codes.Internal, InternalErrorMsg)
			}

			// Кладем в метадату
			md := metadata.Pairs("busy-employees", string(busyEmployeesStr))
			if err := grpc.SetHeader(ctx, md); err != nil {
				a.logger.Error("failed to set header", slog.Any("error", err))
				return nil, status.Error(codes.Internal, InternalErrorMsg)
			}

			// Возвращаем ошибку
			return nil, status.Error(codes.FailedPrecondition, EmployeesAreNotAvailableMsg)
		}

		// Если по такому id нет мероприятия
		if errors.Is(err, eventservice.ErrEventNotFound) {
			a.logger.Info("event not found", slog.Any("eventID", event.ID))
			return nil, status.Error(codes.NotFound, EventNotFoundMsg)
		}

		// Если произошла ошибка
		a.logger.Error("failed to update event", slog.Any("error", err))
		return nil, status.Error(codes.Internal, InternalErrorMsg)
	}

	resp.UpdatedEvent.Room, err = a.getInfoAboutRoom(ctx, resp.UpdatedEvent.Room.ID)
	if err != nil {
		a.logger.Error("failed to get room info", slog.Any("id", resp.UpdatedEvent.Room.ID), slog.Any("error", err))
		return nil, status.Error(codes.Internal, InternalErrorMsg)
	}

	resp.UpdatedEvent.Participants, err = a.getInfoAboutParticipants(ctx, resp.UpdatedEvent.Participants)
	if err != nil {
		a.logger.Error("failed to get participants info", slog.Any("id", resp.UpdatedEvent.Participants), slog.Any("error", err))
		return nil, status.Error(codes.Internal, InternalErrorMsg)
	}

	resp.UpdatedEvent.Creator, err = a.getInfoAboutCreator(ctx, resp.UpdatedEvent.Creator.ID)
	if err != nil {
		a.logger.Error("failed to get creator info", slog.Any("id", resp.UpdatedEvent.Creator.ID), slog.Any("error", err))
		return nil, status.Error(codes.Internal, InternalErrorMsg)
	}

	a.logger.Info("successfully update event", slog.Any("eventID", event.ID))

	return &eventsv1.UpdateResponse{UpdatedEvent: transformEventToResponse(resp.UpdatedEvent)}, nil
}

func (a *api) Delete(ctx context.Context, request *eventsv1.DeleteRequest) (*emptypb.Empty, error) {
	id := request.GetId()

	a.logger.Info("try to delete event by id", slog.Any("id", id))

	if err := a.eventProvider.DeleteByID(ctx, id); err != nil {
		// Если по такому id события нет
		if errors.Is(err, eventservice.ErrEventNotFound) {
			a.logger.Info("event not found", slog.Any("id", id))
			return nil, status.Error(codes.NotFound, EventNotFoundMsg)
		}

		// Произошла ошибка
		a.logger.Error("failed to delete event", slog.Any("id", id), slog.Any("error", err))
		return nil, status.Error(codes.Internal, InternalErrorMsg)
	}

	a.logger.Info("successfully delete event", slog.Any("id", id))

	return &emptypb.Empty{}, nil
}

func transformEventFromRequest(event *eventsv1.Event) *models.Event {
	return &models.Event{
		Title:       event.Title,
		Description: event.Description,
		StartDate:   event.StartDate.AsTime(),
		EndDate:     event.EndDate.AsTime(),
		Room: &models.Room{
			ID: event.Room.Id,
		},
		Creator: &models.Employee{
			ID: event.Creator.Id,
		},
		Participants: func() []*models.Employee {
			participants := make([]*models.Employee, len(event.Participants))
			for i, p := range event.Participants {
				participants[i] = &models.Employee{
					ID: p.Id,
				}
			}
			return participants
		}(),
	}
}

func transformEventToResponse(event *models.Event) *eventsv1.Event {
	return &eventsv1.Event{
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
		Participants: func() []*eventsv1.Employee {
			participants := make([]*eventsv1.Employee, len(event.Participants))
			for i, p := range event.Participants {
				participants[i] = &eventsv1.Employee{
					Id: p.ID,
				}
			}
			return participants
		}(),
	}
}

func (a *api) getInfoAboutRoom(ctx context.Context, roomID uint32) (*models.Room, error) {
	room, err := a.roomsProvider.GetByID(ctx, roomID)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (a *api) getInfoAboutBusyEmployees(ctx context.Context, busyEmployees []uint64) ([]*models.Employee, error) {
	aggregatedBustEmployees, err := a.employeesProvider.GetByIDs(ctx, busyEmployees)
	if err != nil {
		return nil, err
	}

	return aggregatedBustEmployees, nil
}

func (a *api) getInfoAboutParticipants(ctx context.Context, participants []*models.Employee) ([]*models.Employee, error) {
	participantsIDs := make([]uint64, len(participants))
	for i, participant := range participants {
		participantsIDs[i] = participant.ID
	}
	aggregatedParticipants, err := a.employeesProvider.GetByIDs(ctx, participantsIDs)
	if err != nil {
		return nil, err
	}

	return aggregatedParticipants, nil
}

func (a *api) getInfoAboutCreator(ctx context.Context, creatorID uint64) (*models.Employee, error) {
	creator, err := a.employeesProvider.GetByID(ctx, creatorID)
	if err != nil {
		return nil, err
	}

	return creator, nil
}
