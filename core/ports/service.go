package ports

import "nickel/core/domain"

type EntryServicePort interface {
	Create(entry *domain.Entry) (*domain.Entry, error)
	List() ([]domain.Entry, error)
	Delete(id string) error
}
