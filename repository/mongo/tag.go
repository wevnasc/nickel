package mongo

import (
	"nickel/core/domain"
	"nickel/core/errors"
	"nickel/repository/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Tag struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}

type TagList []Tag

func tagFrom(tag *domain.Tag) *Tag {
	ID, _ := primitive.ObjectIDFromHex(tag.ID)

	if ID == primitive.NilObjectID {
		ID = primitive.NewObjectID()
	}

	return &Tag{
		ID:   ID,
		Name: tag.Name,
	}
}

func (s *Tag) toDomain() *domain.Tag {
	return &domain.Tag{
		ID:   s.ID.Hex(),
		Name: s.Name,
	}
}

func tagsFrom(tags []domain.Tag) TagList {
	list := make(TagList, len(tags))

	for idx, t := range tags {

		list[idx] = *tagFrom(&t)
	}

	return list
}

func (s TagList) toDomain() []domain.Tag {
	list := make([]domain.Tag, len(s))
	for idx, t := range s {
		tag := t.toDomain()
		list[idx] = *tag
	}
	return list
}

type TagRepository struct {
	timeout time.Duration
	client  *mongo.Client
	coll    *mongo.Collection
}

func NewTagRepository(client *mongo.Client, database string, timeout time.Duration) *TagRepository {
	return &TagRepository{
		client:  client,
		timeout: timeout,
		coll:    client.Database(database).Collection("tags"),
	}
}

func (m *TagRepository) CreateMany(tags []domain.Tag) ([]domain.Tag, error) {
	ctx, cancel := config.TimeoutContext(m.timeout)
	defer cancel()

	schemaList := tagsFrom(tags)

	var insertList []interface{}

	for _, s := range schemaList {
		insertList = append(insertList, s)
	}

	data, err := m.coll.InsertMany(ctx, insertList)

	if err != nil {
		return nil, errors.Wrap(errors.InsertData, "not was possible to create the tags", err)
	}

	for idx, t := range schemaList {
		if oid, ok := data.InsertedIDs[idx].(primitive.ObjectID); ok {
			t.ID = oid
		}
	}

	return schemaList.toDomain(), nil
}

func (m *TagRepository) GetByNames(names []string) ([]domain.Tag, error) {
	ctx, cancel := config.TimeoutContext(m.timeout)
	defer cancel()

	cursor, err := m.coll.Find(ctx, bson.M{"name": bson.M{"$in": names}})

	if cursor.Err() != nil {
		return nil, errors.Wrap(errors.FindData, "not was possible to query data", cursor.Err())
	}

	var list TagList

	for cursor.Next(ctx) {
		var tag Tag
		err := cursor.Decode(&tag)
		if err != nil {
			return list.toDomain(), errors.Wrap(errors.Serialization, "not was to decode the tag from database", err)
		}
		list = append(list, tag)
	}

	if err != nil {
		return nil, errors.Wrap(errors.Serialization, "not was possible decode db tag", err)
	}

	return list.toDomain(), nil
}
