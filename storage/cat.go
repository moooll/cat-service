package storage

import (
	"errors"
	"sync"

	"cat-crud/models"

	"github.com/google/uuid"
)

type Catalog struct {
	str map[uuid.UUID]models.Cat
	mu  *sync.Mutex
}

func NewCatalog() *Catalog {
	str := make(map[uuid.UUID]models.Cat)
	return &Catalog{
		str,
		&sync.Mutex{},
	}
}

func (s *Catalog) Save(cat models.Cat) {
	id := uuid.New()
	s.mu.Lock()
	s.str[id] = cat
	s.mu.Unlock()
}

func (s *Catalog) Delete(id uuid.UUID) {
	s.mu.Lock()
	delete(s.str, id)
	s.mu.Unlock()
}

func (s *Catalog) Update(id uuid.UUID, fields []interface{}) error {
	cat, ok := s.str[id]
	if !ok {
		return errors.New("No such item in the storage!")
	}
	type field struct {
		name  string
		value string
	}
	cat.
	
}
