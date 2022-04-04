package tag

import (
	"nickel/core/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TagSchema struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}

type TagListSchema []TagSchema

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
