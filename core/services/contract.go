package services

import "nickel/core/domain"

type EntryService interface {
	Create(entry *domain.Entry) (*domain.Entry, error)
	List() ([]domain.Entry, error)
	Delete(id string) error
}
