package mongo

import (
	"nickel/core/domain"
	"nickel/core/errors"
	"nickel/repositories/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoEntryRepository struct {
	timeout time.Duration
	client  *mongo.Client
	coll    *mongo.Collection
}

func NewMongoEntryRepository(client *mongo.Client, database string, timeout time.Duration) *MongoEntryRepository {
	return &MongoEntryRepository{
		client:  client,
		timeout: timeout,
		coll:    client.Database(database).Collection("entries"),
	}
}

func (m *MongoEntryRepository) Create(entry *domain.Entry) (*domain.Entry, error) {
	ctx, cancel := config.TimeoutContext(m.timeout)
	defer cancel()

	entrySchema := entryFrom(entry)
	data, err := m.coll.InsertOne(ctx, entrySchema)

	if err != nil {
		return nil, errors.Wrap(errors.InsertData, "not was possible to create the entry", err)
	}

	if oid, ok := data.InsertedID.(primitive.ObjectID); ok {
		entrySchema.ID = oid
	}

	return entrySchema.toDomain(), nil
}

func (m *MongoEntryRepository) List() ([]domain.Entry, error) {
	ctx, cancel := config.TimeoutContext(m.timeout)
	defer cancel()

	cursor, err := m.coll.Find(ctx, bson.D{})

	if err != nil {
		return nil, errors.Wrap(errors.FindData, "not was possible to query the entries", err)
	}

	var list EntryListSchema

	for cursor.Next(ctx) {
		var entry EntrySchema
		err := cursor.Decode(&entry)
		if err != nil {
			return list.toDomain(), errors.Wrap(errors.Serialization, "not was to decode the entry from database", err)
		}
		list = append(list, entry)
	}

	return list.toDomain(), nil
}

func (m *MongoEntryRepository) Get(id string) (*domain.Entry, error) {
	ctx, cancel := config.TimeoutContext(m.timeout)
	defer cancel()

	ID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, errors.Wrap(errors.InvalidIdentity, "the ID is invalid or is malformed", err)
	}

	res := m.coll.FindOne(ctx, bson.M{"_id": ID})

	if res.Err() != nil {
		return nil, errors.Wrap(errors.NotFound, "not was possible to found entry", res.Err())
	}

	entry := &EntrySchema{}
	err = res.Decode(entry)

	if err != nil {
		return nil, errors.Wrap(errors.Serialization, "not was possible decode db entry", err)
	}

	return entry.toDomain(), nil
}

func (m *MongoEntryRepository) Delete(id string) error {
	ctx, cancel := config.TimeoutContext(m.timeout)
	defer cancel()

	ID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return errors.Wrap(errors.InvalidIdentity, "the ID is invalid or is malformed", err)
	}

	_, err = m.coll.DeleteOne(ctx, bson.M{"_id": ID})

	if err != nil {
		return errors.Wrap(errors.DeleteData, "not was possible to delete entry", err)
	}

	return nil
}
