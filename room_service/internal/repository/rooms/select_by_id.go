package rooms

import (
	"context"
	"errors"
	"fmt"
	"github.com/Sleeps17/events-planning-service-backend/rooms_service/internal/domain/models"
	repo "github.com/Sleeps17/events-planning-service-backend/rooms_service/internal/repository"
	"github.com/jackc/pgx/v5"
)

func (r *repository) SelectByID(ctx context.Context, id uint64) (*models.Room, error) {
	const op = "postgresql.Room"

	query := `SELECT id, room_name, capacity FROM rooms_schema.rooms WHERE id = $1`

	var room models.Room
	row := r.pool.QueryRow(ctx, query, id)

	if err := row.Scan(&room.ID, &room.Name, &room.Capacity); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repo.ErrRoomNotFound
		}

		return nil, fmt.Errorf("%s: can't select room: %w", op, err)
	}

	return &room, nil
}
