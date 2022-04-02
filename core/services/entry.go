package services

import (
	"fmt"
	"nickel/core/domain"
	"nickel/core/errors"
	"nickel/core/ports"
)

type EntryServiceAdapter struct {
	repo ports.EntryRepositoryPort
}

func NewEntryService(repo ports.EntryRepositoryPort) ports.EntryServicePort {
	return &EntryServiceAdapter{repo: repo}
}

func (s *EntryServiceAdapter) Create(entry *domain.Entry) (*domain.Entry, error) {
	return s.repo.Create(entry)
}

func (s *EntryServiceAdapter) List() ([]domain.Entry, error) {
	return s.repo.List()
}

func (s EntryServiceAdapter) Delete(id string) error {

	_, err := s.repo.Get(id)

	if err != nil {
		return errors.New(errors.NotFound, fmt.Sprintf("the entry with id %s was not found", id))
	}

	return s.repo.Delete(id)
}
