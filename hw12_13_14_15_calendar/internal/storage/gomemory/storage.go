package gomemory

import (
	"fmt"
	"sync"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	models "hw12_13_14_15_calendar/internal/models"
)

type inMemoryDatabase map[int]models.Event

type Storage struct {
	data inMemoryDatabase
	mu   *sync.RWMutex
}

func NewStorage(_ string) interfaces.Storager {
	return &Storage{nil, &sync.RWMutex{}}
}

func (s *Storage) Connect() error {
	s.data = make(inMemoryDatabase)
	return nil
}

func (s *Storage) Close() error {
	s.data = nil
	return nil
}

func (s *Storage) CreateEvent(e *models.Event) error {
	if s.data == nil {
		return fmt.Errorf("may by is need to reconnect")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	e.PK = len(s.data) + 1
	s.data[e.PK] = *e
	return nil // TODO: not forever nil for may be UNIQUE precheck
}

func (s *Storage) ReadEvent(pk int) (*models.Event, error) {
	if s.data == nil {
		return nil, fmt.Errorf("may by is need to reconnect")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	event, exists := s.data[pk]
	if exists {
		return &event, nil
	}
	// return nil, fmt.Errorf("it does not exists")
	return nil, nil
}

func (s *Storage) UpdateEvent(e *models.Event) error {
	if s.data == nil {
		return fmt.Errorf("may by is need to reconnect")
	}
	if e.PK == 0 {
		return fmt.Errorf("it is not idented")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[e.PK] = *e
	return nil
}

func (s *Storage) DeleteEvent(e *models.Event) error {
	if s.data == nil {
		return fmt.Errorf("may by is need to reconnect")
	}
	if e.PK == 0 {
		return fmt.Errorf("it is not idented")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, e.PK)
	return nil
}

func (s *Storage) ListEvents() ([]models.Event, error) {
	if s.data == nil {
		return nil, fmt.Errorf("may by is need to reconnect")
	}
	events := []models.Event{}
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, event := range s.data {
		events = append(events, event)
	}
	return events, nil // TODO: no events without connect
}
