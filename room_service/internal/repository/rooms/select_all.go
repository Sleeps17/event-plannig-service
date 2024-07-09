package rooms

import (
	"context"
	"errors"
	"fmt"
	"github.com/Sleeps17/events-planning-service-backend/rooms_service/internal/domain/models"
	"github.com/jackc/pgx/v5"
)

func (r *repository) SelectAll(ctx context.Context) ([]*models.Room, error) {
	query := `SELECT id, room_name, capacity FROM rooms_schema.rooms`

	rooms := make([]*models.Room, 0)
	rows, err := r.pool.Query(ctx, query)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*models.Room{}, nil
		}
		return nil, fmt.Errorf("failed to select rooms: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var room models.Room
		if err := rows.Scan(&room.ID, &room.Name, &room.Capacity); err != nil {
			return nil, fmt.Errorf("failed to select rooms: %w", err)
		}

		rooms = append(rooms, &room)
	}

	return rooms, nil
}
