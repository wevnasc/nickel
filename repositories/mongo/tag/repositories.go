package tag

import (
	"nickel/core/domain"
	"nickel/core/errors"
	"nickel/repositories/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoTagRepository struct {
	timeout time.Duration
	client  *mongo.Client
	coll    *mongo.Collection
}

func NewMongoTagRepository(client *mongo.Client, database string, timeout time.Duration) *MongoTagRepository {
	return &MongoTagRepository{
		client:  client,
		timeout: timeout,
		coll:    client.Database(database).Collection("tags"),
	}
}

func (m *MongoTagRepository) CreateMany(tags []domain.Tag) ([]domain.Tag, error) {
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

func (m *MongoTagRepository) GetByNames(names []string) ([]domain.Tag, error) {
	ctx, cancel := config.TimeoutContext(m.timeout)
	defer cancel()

	cursor, err := m.coll.Find(ctx, bson.M{"name": bson.M{"$in": names}})

	if cursor.Err() != nil {
		return nil, errors.Wrap(errors.FindData, "not was possible to query data", cursor.Err())
	}

	var list TagListSchema

	for cursor.Next(ctx) {
		var tag TagSchema
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
