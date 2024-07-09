package api

import (
	"context"
	"errors"
	"github.com/Sleeps17/events-planning-service-backend/rooms_service/internal/domain/models"
	roomservice "github.com/Sleeps17/events-planning-service-backend/rooms_service/internal/services/rooms"
	roomsv1 "github.com/Sleeps17/events-planning-service-backend/rooms_service/protos/gen/go/rooms"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
)

type api struct {
	roomsv1.UnimplementedRoomsServer
	logger        *slog.Logger
	roomsProvider roomservice.Service
}

func Register(
	srv *grpc.Server,
	logger *slog.Logger,
	roomsProvider roomservice.Service,
) {
	roomsv1.RegisterRoomsServer(srv, &api{
		logger:        logger,
		roomsProvider: roomsProvider,
	})
}

func (a *api) Create(ctx context.Context, request *roomsv1.CreateRequest) (*roomsv1.CreateResponse, error) {
	room := transformRoomFromRequest(request.GetRoom())

	a.logger.Info("try to handle create request", slog.Any("room", room))

	resp, err := a.roomsProvider.Add(ctx, room)

	if err != nil {
		if errors.Is(err, roomservice.ErrRoomExists) {
			a.logger.Error("room already exists", slog.Any("room", room))
			return nil, status.Error(codes.AlreadyExists, RoomAlreadyExistsMsg)
		}

		a.logger.Error("failed to add room", slog.Any("error", err))
		return nil, status.Error(codes.Internal, InternalErrorMsg)
	}

	a.logger.Info("successfully added room", slog.Any("roomID", room.ID))
	return &roomsv1.CreateResponse{Id: resp}, nil
}

func (a *api) GetByID(ctx context.Context, request *roomsv1.GetRequest) (*roomsv1.GetResponse, error) {
	id := request.GetId()

	a.logger.Info("try to handle get room by id", slog.Any("id", id))

	room, err := a.roomsProvider.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, roomservice.ErrRoomNotFound) {
			a.logger.Error("room not found", slog.Any("id", id))
			return nil, status.Error(codes.NotFound, RoomNotFoundMsg)
		}

		a.logger.Error("failed to get room by id", slog.Any("id", id))
		return nil, status.Error(codes.Internal, InternalErrorMsg)
	}

	a.logger.Info("successfully got room by id", slog.Any("room", room))
	return &roomsv1.GetResponse{Room: transformRoomToRequest(room)}, nil
}

func (a *api) GetAll(ctx context.Context, _ *emptypb.Empty) (*roomsv1.GetAllResponse, error) {
	a.logger.Info("try to handle get all rooms")

	rooms, err := a.roomsProvider.GetAll(ctx)
	if err != nil {
		a.logger.Error("failed to get all rooms", slog.Any("error", err))
		return nil, status.Error(codes.Internal, InternalErrorMsg)
	}

	a.logger.Info("successfully got all rooms", slog.Any("rooms", rooms))
	return &roomsv1.GetAllResponse{Rooms: transformRoomsToRequest(rooms)}, nil
}

func (a *api) Update(ctx context.Context, request *roomsv1.UpdateRequest) (*roomsv1.UpdateResponse, error) {
	id := request.GetId()
	updatedRoom := transformRoomFromRequest(request.GetUpdatedRoom())

	a.logger.Info("try to handle update room", slog.Any("id", id))

	err := a.roomsProvider.Update(ctx, updatedRoom)
	if err != nil {
		a.logger.Error("failed to update room", slog.Any("id", id))
		return nil, status.Error(codes.Internal, InternalErrorMsg)
	}

	a.logger.Info("successfully updated room", slog.Any("id", id))
	return &roomsv1.UpdateResponse{UpdatedRoom: transformRoomToRequest(updatedRoom)}, nil
}

func (a *api) Delete(ctx context.Context, request *roomsv1.DeleteRequest) (*emptypb.Empty, error) {
	id := request.GetId()

	a.logger.Info("try to handle delete room", slog.Any("id", id))

	err := a.roomsProvider.DeleteByID(ctx, id)
	if err != nil {
		a.logger.Error("failed to delete room", slog.Any("id", id))
		return nil, status.Error(codes.Internal, InternalErrorMsg)
	}

	a.logger.Info("successfully deleted room", slog.Any("id", id))
	return &emptypb.Empty{}, nil
}

func transformRoomFromRequest(room *roomsv1.Room) *models.Room {
	return &models.Room{
		Name:     room.Name,
		Capacity: room.Capacity,
	}
}

func transformRoomToRequest(room *models.Room) *roomsv1.Room {
	return &roomsv1.Room{
		Name:     room.Name,
		Capacity: room.Capacity,
	}
}

func transformRoomsToRequest(rooms []*models.Room) []*roomsv1.Room {
	roomResponses := make([]*roomsv1.Room, len(rooms))
	for i, room := range rooms {
		roomResponses[i] = transformRoomToRequest(room)
	}

	return roomResponses
}
