package services

import (
	"nickel/core/domain"
	"nickel/core/ports"
)

type EntryServiceAdapter struct {
	repo ports.EntryRepositoryPort
}

func NewEntryService(repo ports.EntryRepositoryPort) ports.EntryServicePort {
	return &EntryServiceAdapter{repo: repo}
}

func (s *EntryServiceAdapter) Create(entry *domain.Entry) (*domain.Entry, error) {
	entry, err := s.repo.Create(entry)

	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (s *EntryServiceAdapter) List() ([]domain.Entry, error) {
	entries, err := s.repo.List()

	if err != nil {
		return nil, err
	}

	return entries, nil
}
