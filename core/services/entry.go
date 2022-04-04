package services

import (
	"fmt"
	"nickel/core/domain"
	"nickel/core/errors"
	"nickel/core/ports"
)

type EntryServiceAdapter struct {
	entryRepo ports.EntryRepositoryPort
	tagRepo   ports.TagRepositoryPort
}

func NewEntryService(entryRepo ports.EntryRepositoryPort, tagRepo ports.TagRepositoryPort) ports.EntryServicePort {
	return &EntryServiceAdapter{
		entryRepo: entryRepo,
		tagRepo:   tagRepo,
	}
}

func (s *EntryServiceAdapter) findOrCreateTags(tags []domain.Tag) ([]domain.Tag, error) {
	names := domain.Names(tags)
	tagsFound, err := s.tagRepo.GetByNames(names)

	if err != nil {
		return nil, err
	}

	missing := domain.MissingTags(tagsFound, tags)
	var newTags []domain.Tag

	if len(missing) > 0 {
		newTags, err = s.tagRepo.CreateMany(missing)
	}

	if err != nil {
		return nil, err
	}

	return append(tagsFound, newTags...), nil
}

func (s *EntryServiceAdapter) Create(entry *domain.Entry) (*domain.Entry, error) {
	tags, err := s.findOrCreateTags(entry.Tags)

	if err != nil {
		return nil, err
	}

	entry.Tags = tags
	return s.entryRepo.Create(entry)
}

func (s *EntryServiceAdapter) List() ([]domain.Entry, error) {
	return s.entryRepo.List()
}

func (s EntryServiceAdapter) Delete(id string) error {

	_, err := s.entryRepo.Get(id)

	if err != nil {
		return errors.New(errors.NotFound, fmt.Sprintf("the entry with id %s was not found", id))
	}

	return s.entryRepo.Delete(id)
}
