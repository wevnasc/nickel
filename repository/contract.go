package repository

import "nickel/core/domain"

type Entry interface {
	Create(entry *domain.Entry) (*domain.Entry, error)
	List() ([]domain.Entry, error)
	Get(id string) (*domain.Entry, error)
	Delete(id string) error
}

type Tag interface {
	CreateMany(tag []domain.Tag) ([]domain.Tag, error)
	GetByNames(names []string) ([]domain.Tag, error)
}
