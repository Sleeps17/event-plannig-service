package postgresql_events_repository

var (
	SelectEventByIDQuery = `
        SELECT e.id, e.title, e.description, e.start_date, e.end_date, 
            e.room_id, e.creator_id, 
            ARRAY_AGG(p.employee_id) AS participants
        FROM events_schema.events e
        LEFT JOIN events_schema.participants p ON e.id = p.event_id
        WHERE e.id = $1
        GROUP BY e.id
    `
	SelectEventsBetweenTwoDatesQuery = `SELECT e.id, e.title, e.description, e.start_date, e.end_date,
										e.room_id, e.creator_id,
										ARRAY_AGG(p.employee_id) AS participants
										FROM events_schema.events e
										LEFT JOIN events_schema.participants p ON e.id = p.event_id
										WHERE e.start_date BETWEEN $1 AND $2;
`
	CheckRoomIsAvailableQuery = `SELECT NOT EXISTS (
    							 	SELECT 1
    								FROM events_schema.events e
    								WHERE e.room_id = $1
    								AND (
        								($2 < e.end_date) AND
        								($3 > e.start_date)
    								)
								);
`
	CheckEmployeesAreAvailableQuery = `SELECT DISTINCT employee_id
										FROM events_schema.events e
										JOIN events_schema.participants p ON e.id = p.event_id
										WHERE p.employee_id = ANY($1)
    									AND (
        									(e.start_date < $2 AND e.end_date > $2)
        									OR (e.start_date < $3 AND e.end_date > $3)
        									OR (e.start_date >= $2 AND e.end_date <= $3)
    									);`
	InsertEventQuery             = `INSERT INTO events_schema.events(title, description, start_date, end_date, room_id, creator_id) VALUES($1, $2, $3, $4, $5, $6) RETURNING id`
	InsertEventParticipantsQuery = `INSERT INTO events_schema.participants(event_id, employee_id) SELECT $1, unnest($2::INT[]);`
	UpdateEventQuery             = `UPDATE events_schema.events SET title = $1, description = $2, start_date = $3, end_date = $4, room_id = $5, creator_id = $6 WHERE id = $7;`
	DeleteEventQuery             = `DELETE FROM events_schema.events WHERE id = $1;`
	DeleteEventParticipantsQuery = `DELETE FROM events_schema.participants WHERE event_id = $1;`
)
