package rooms

import (
	"context"
	"errors"
	"fmt"
	"github.com/Sleeps17/events-planning-service-backend/rooms_service/internal/domain/models"
	repo "github.com/Sleeps17/events-planning-service-backend/rooms_service/internal/repository"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	uniqueViolationCode = "23505"

	defaultIDValue = 0
)

func (r *repository) Insert(ctx context.Context, room *models.Room) (uint64, error) {
	const op = "room-repository.Save"

	query := `INSERT INTO rooms_schema.rooms(room_name, capacity) VALUES($1, $2) RETURNING id`

	var roomId uint64
	row := r.pool.QueryRow(ctx, query, room.Name, room.Capacity)

	if err := row.Scan(&roomId); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
			return defaultIDValue, repo.ErrRoomExists
		}

		return defaultIDValue, fmt.Errorf("%s: can't add room: %w", op, err)
	}

	return roomId, nil
}
