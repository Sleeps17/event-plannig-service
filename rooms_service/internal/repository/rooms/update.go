package rooms

import (
	"context"
	"errors"
	"fmt"
	"github.com/Sleeps17/events-planning-service-backend/rooms_service/internal/domain/models"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *repository) Update(ctx context.Context, room *models.Room) error {
	const op = "room-repository.Update"

	query := `UPDATE rooms_schema.rooms SET room_name = $2, capacity = $3 WHERE id = $1`

	if _, err := r.pool.Exec(ctx, query, room.ID, room.Name, room.Capacity); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return fmt.Errorf("%s: database error: %w", op, pgErr)
		}

		return fmt.Errorf("%s: can't update room: %w", op, err)
	}

	return nil
}
