package postgresql_events_repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/Sleeps17/events-planning-service-backend/events_service/internal/domain/models"
	eventsrepository "github.com/Sleeps17/events-planning-service-backend/events_service/internal/repository/events"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

const (
	UniqueConstraintErrorCode = "23505"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) SelectByID(ctx context.Context, id uint64) (*models.Event, error) {
	var event models.Event
	var participants []uint64
	if err := r.db.QueryRow(ctx, SelectEventByIDQuery, id).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.StartDate,
		&event.EndDate,
		&event.Room.ID,
		&event.Creator.ID,
		&participants,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, eventsrepository.ErrEventNotFound
		}
		return nil, fmt.Errorf("failed to select event by id: %w", err)
	}
	for _, id := range participants {
		event.Participants = append(event.Participants, &models.Employee{ID: id})
	}

	return &event, nil
}

func (r *Repository) SelectAllBetweenTwoDates(ctx context.Context, startDate, endDate time.Time) ([]*models.Event, error) {
	rows, err := r.db.Query(ctx, SelectEventsBetweenTwoDatesQuery, startDate, endDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*models.Event{}, nil
		}
		return nil, fmt.Errorf("failed to select events between two dates: %w", err)
	}

	events := make([]*models.Event, 0)
	for rows.Next() {
		var event models.Event
		var participants []uint64
		if err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.StartDate,
			&event.EndDate,
			&event.Room.ID,
			&event.Creator.ID,
			&participants,
		); err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		for _, id := range participants {
			event.Participants = append(event.Participants, &models.Employee{ID: id})
		}

		events = append(events, &event)
	}

	return events, nil
}

func (r *Repository) CheckRoomIsAvailable(ctx context.Context, roomID uint32, startDate, endDate time.Time) (bool, error) {
	var roomIsAvailable bool
	if err := r.db.QueryRow(ctx, CheckRoomIsAvailableQuery, roomID, startDate, endDate).Scan(&roomIsAvailable); err != nil {
		return false, fmt.Errorf("failed to check room availability: %w", err)
	}
	return roomIsAvailable, nil
}

func (r *Repository) CheckEmployeesAreAvailable(ctx context.Context, employeeIDs []uint64, startDate, endDate time.Time) ([]uint64, error) {
	rows, err := r.db.Query(ctx, CheckEmployeesAreAvailableQuery, employeeIDs, startDate, endDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []uint64{}, nil
		}
		return nil, fmt.Errorf("failed to check employees availability: %w", err)
	}
	defer rows.Close()

	var busyEmployeeIDs []uint64
	for rows.Next() {
		var employeeID uint64
		if err := rows.Scan(&employeeID); err != nil {
			return nil, fmt.Errorf("failed to scan employee ID: %w", err)
		}
		busyEmployeeIDs = append(busyEmployeeIDs, employeeID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over result rows: %w", err)
	}

	return busyEmployeeIDs, nil
}

func (r *Repository) Insert(ctx context.Context, event *models.Event) (id uint64, err error) {
	var tx pgx.Tx
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	tx, err = r.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	if err := tx.QueryRow(
		ctx,
		InsertEventQuery,
		event.Title,
		event.Description,
		event.StartDate,
		event.EndDate,
		event.Room.ID,
		event.Creator.ID,
	).Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to insert event: %w", err)
	}

	participants := make([]uint64, len(event.Participants))
	for i, p := range event.Participants {
		participants[i] = p.ID
	}

	_, err = tx.Exec(ctx, InsertEventParticipantsQuery, id, participants)
	if err != nil {
		return 0, fmt.Errorf("failed to insert event participants: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return id, nil
}

func (r *Repository) UpdateByID(ctx context.Context, updatedEvent *models.Event) (err error) {
	var tx pgx.Tx
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	tx, err = r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	cmd, err := tx.Exec(ctx, UpdateEventQuery,
		updatedEvent.Title,
		updatedEvent.Description,
		updatedEvent.StartDate,
		updatedEvent.EndDate,
		updatedEvent.Room.ID,
		updatedEvent.Creator.ID,
		updatedEvent.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	} else if cmd.RowsAffected() == 0 {
		return eventsrepository.ErrEventNotFound
	}

	_, err = tx.Exec(ctx, DeleteEventParticipantsQuery, updatedEvent.ID)
	if err != nil {
		return fmt.Errorf("failed to delete event participants: %w", err)
	}

	participants := make([]uint64, len(updatedEvent.Participants))
	for i, p := range updatedEvent.Participants {
		participants[i] = p.ID
	}

	_, err = tx.Exec(ctx, InsertEventParticipantsQuery, updatedEvent.ID, participants)
	if err != nil {
		return fmt.Errorf("failed to insert event participants: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *Repository) DeleteByID(ctx context.Context, id uint64) (err error) {
	var tx pgx.Tx
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	tx, err = r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	_, err = tx.Exec(ctx, DeleteEventParticipantsQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete event participants: %w", err)
	}

	cmd, err := tx.Exec(ctx, DeleteEventQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return eventsrepository.ErrEventNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
