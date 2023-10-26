package pgsqldtb

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq" // a blank import should be justifying.
	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	models "hw12_13_14_15_calendar/internal/models"
)

type Storage struct { // TODO
	dsn        string
	connection *sql.DB
}

func NewStorage(dsn string) interfaces.Storager {
	return &Storage{dsn, nil}
}

func (s *Storage) Connect() error {
	db, err := sql.Open("postgres", s.dsn)
	if err != nil {
		return err
	}
	s.connection = db
	return nil
}

func (s *Storage) Close() error {
	if s.connection == nil {
		return nil
	}
	err := s.connection.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) CreateEvent(e *models.Event) (*models.Event, error) {
	sqlStatement := `
	INSERT INTO hw15calendar.events(
		"title", "description", "startat", "durationseconds", "owner", "notifyearlyseconds", "sheduled"
	) values($1, $2, $3, $4, $5, $6, $7)
	RETURNING "pk";`
	err := s.connection.QueryRow(sqlStatement,
		e.Title, e.Description, e.StartAt, e.Duration, e.Owner, e.NotifyEarly, &e.Sheduled,
	).Scan(&(e.PK))
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (s *Storage) ReadEvent(pk int) (*models.Event, error) {
	var e models.Event
	sqlStatement := `
	SELECT "pk", "title", "description", "startat", "durationseconds", "owner", "notifyearlyseconds", "sheduled"
	FROM hw15calendar.events WHERE "pk"=$1;`
	row := s.connection.QueryRow(sqlStatement, pk)
	err := row.Scan(&(e.PK),
		&(e.Title), &(e.Description), &(e.StartAt), &(e.Duration),
		&(e.Owner), &(e.NotifyEarly), &(e.Sheduled))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

func (s *Storage) UpdateEvent(e *models.Event) (*models.Event, error) {
	if e.PK == 0 {
		return nil, fmt.Errorf("it is not idented")
	}
	sqlStatement := `
	UPDATE hw15calendar.events 
	SET 
	"title"=$1, "description"=$2, "startat"=$3, "durationseconds"=$4, 
	"owner"=$5, "notifyearlyseconds"=$6, "sheduled"=$7
	WHERE pk=$8;`
	_, err := s.connection.Exec(sqlStatement, e.Title, e.Description,
		e.StartAt, e.Duration, e.Owner, e.NotifyEarly, e.Sheduled,
		e.PK)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (s *Storage) DeleteEvent(e *models.Event) (*models.Event, error) {
	if e.PK == 0 {
		return nil, fmt.Errorf("it is not idented")
	}
	sqlStatement := `DELETE FROM hw15calendar.events WHERE pk=$1;`
	_, err := s.connection.Exec(sqlStatement, e.PK)
	if err != nil {
		return nil, err
	}
	e.PK = 0
	return e, nil
}

func (s *Storage) ListEvents() ([]models.Event, error) {
	var e models.Event
	var events []models.Event
	sqlStatement := `
	SELECT "pk", "title", "description", "startat", "durationseconds", "owner", "notifyearlyseconds", "sheduled"
	FROM hw15calendar.events;`
	rows, err := s.connection.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&e.PK, &e.Title, &e.Description, &e.StartAt, &e.Duration, &e.Owner, &e.NotifyEarly, &e.Sheduled)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

func (s *Storage) ListNotSheduledEvents() ([]models.Event, error) {
	var e models.Event
	var events []models.Event
	sqlStatement := `
	SELECT "pk", "title", "description", "startat", "durationseconds", "owner", "notifyearlyseconds", "sheduled"
	FROM hw15calendar.events
	WHERE "sheduled" IS NOT TRUE;`
	rows, err := s.connection.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&e.PK, &e.Title, &e.Description, &e.StartAt, &e.Duration, &e.Owner, &e.NotifyEarly, &e.Sheduled)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}
