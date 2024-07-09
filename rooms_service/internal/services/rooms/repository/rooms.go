package repository

import (
	"context"
	"errors"
	"github.com/Sleeps17/events-planning-service-backend/rooms_service/internal/domain/models"
	roomsrepository "github.com/Sleeps17/events-planning-service-backend/rooms_service/internal/repository"
	roomsservice "github.com/Sleeps17/events-planning-service-backend/rooms_service/internal/services/rooms"
)

type Service struct {
	repo roomsrepository.Repository
}

func New(repo roomsrepository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Add(ctx context.Context, room *models.Room) (uint64, error) {
	id, err := s.repo.Insert(ctx, room)

	if err != nil {
		if errors.Is(err, roomsrepository.ErrRoomExists) {
			return 0, roomsservice.ErrRoomExists
		}

		return 0, err
	}

	return id, nil
}

func (s *Service) GetByID(ctx context.Context, id uint64) (*models.Room, error) {
	room, err := s.repo.SelectByID(ctx, id)

	if err != nil {
		if errors.Is(err, roomsrepository.ErrRoomNotFound) {
			return nil, roomsservice.ErrRoomNotFound
		}

		return nil, err
	}

	return room, nil
}

func (s *Service) GetAll(ctx context.Context) ([]*models.Room, error) {
	rooms, err := s.repo.SelectAll(ctx)

	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (s *Service) Update(ctx context.Context, updatedRoom *models.Room) error {
	err := s.repo.Update(ctx, updatedRoom)

	return err
}

func (s *Service) DeleteByID(ctx context.Context, id uint64) error {
	err := s.repo.Delete(ctx, id)

	return err
}
