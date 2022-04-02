package mongo

import (
	"nickel/core/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TagSchema struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}

type TagListSchema []TagSchema

type EntrySchema struct {
	ID          primitive.ObjectID `bson:"_id"`
	Description string             `bson:"description"`
	Amount      float64            `bson:"amount"`
	Tags        TagListSchema      `bson:"tags"`
	Type        string             `bson:"type"`
}

type EntryListSchema []EntrySchema

func entryFrom(entry *domain.Entry) *EntrySchema {
	ID, _ := primitive.ObjectIDFromHex(entry.ID)

	if ID == primitive.NilObjectID {
		ID = primitive.NewObjectID()
	}

	return &EntrySchema{
		ID:          ID,
		Description: entry.Description,
		Amount:      entry.Amount,
		Tags:        tagsFrom(entry.Tags),
		Type:        string(entry.Type),
	}

}

func (s *EntrySchema) toDomain() *domain.Entry {
	return &domain.Entry{
		ID:          s.ID.Hex(),
		Description: s.Description,
		Amount:      s.Amount,
		Tags:        s.Tags.toDomain(),
		Type:        domain.Type(s.Type),
	}
}

func (s EntryListSchema) toDomain() []domain.Entry {
	list := make([]domain.Entry, len(s))
	for idx, e := range s {
		entry := e.toDomain()
		list[idx] = *entry
	}

	return list
}

func tagFrom(tag *domain.Tag) *TagSchema {
	ID, _ := primitive.ObjectIDFromHex(tag.ID)

	if ID == primitive.NilObjectID {
		ID = primitive.NewObjectID()
	}

	return &TagSchema{
		ID:   ID,
		Name: tag.Name,
	}
}

func (s *TagSchema) toDomain() *domain.Tag {
	return &domain.Tag{
		ID:   s.ID.Hex(),
		Name: s.Name,
	}
}

func tagsFrom(tags []domain.Tag) TagListSchema {
	list := make(TagListSchema, len(tags))

	for idx, t := range tags {

		list[idx] = *tagFrom(&t)
	}

	return list
}

func (s TagListSchema) toDomain() []domain.Tag {
	list := make([]domain.Tag, len(s))
	for idx, t := range s {
		tag := t.toDomain()
		list[idx] = *tag
	}
	return list
}
