package gomemory

import (
	"fmt"
	"sync"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	models "hw12_13_14_15_calendar/internal/models"
)

type inMemoryDatabase map[int]*models.Event

type Storage struct {
	data     inMemoryDatabase
	mu       *sync.RWMutex
	sequence int
}

func NewStorage() interfaces.Storager {
	return &Storage{nil, &sync.RWMutex{}, 0}
}

func (s *Storage) Connect() error {
	if s.data == nil {
		s.data = make(inMemoryDatabase)
	}
	return nil
}

func (s *Storage) Close() error {
	return nil
}

func (s *Storage) CreateEvent(e *models.Event) (*models.Event, error) {
	if s.data == nil {
		return nil, fmt.Errorf("may by is need to reconnect")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sequence++
	e.PK = s.sequence
	s.data[e.PK] = e
	return e, nil // TODO: not forever nil for may be UNIQUE precheck
}

func (s *Storage) ReadEvent(pk int) (*models.Event, error) {
	if s.data == nil {
		return nil, fmt.Errorf("may by is need to reconnect")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	event, exists := s.data[pk]
	if exists {
		return event, nil
	}
	return nil, nil
}

func (s *Storage) UpdateEvent(e *models.Event) (*models.Event, error) {
	if s.data == nil {
		return nil, fmt.Errorf("may by is need to reconnect")
	}
	if e.PK == 0 {
		return nil, fmt.Errorf("it is not idented")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[e.PK] = e
	return e, nil
}

func (s *Storage) DeleteEvent(e *models.Event) (*models.Event, error) {
	if s.data == nil {
		return nil, fmt.Errorf("may by is need to reconnect")
	}
	if e.PK == 0 {
		return nil, fmt.Errorf("it is not idented")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	e = s.data[e.PK]
	delete(s.data, e.PK)
	return e, nil
}

func (s *Storage) ListEvents() ([]models.Event, error) {
	if s.data == nil {
		return nil, fmt.Errorf("may by is need to reconnect")
	}
	events := []models.Event{}
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, event := range s.data {
		events = append(events, *event)
	}
	return events, nil
}
