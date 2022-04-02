package ports

import "nickel/core/domain"

type EntryRepositoryPort interface {
	Create(entry *domain.Entry) (*domain.Entry, error)
	List() ([]domain.Entry, error)
	Get(id string) (*domain.Entry, error)
	Delete(id string) error
}
