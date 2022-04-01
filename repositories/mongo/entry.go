package mongo

import (
	"context"
	"nickel/core/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		ID, _ := primitive.ObjectIDFromHex(t.ID)

		if ID == primitive.NilObjectID {
			ID = primitive.NewObjectID()
		}

		list[idx] = TagSchema{
			ID:   ID,
			Name: t.Name,
		}
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

func timeoutContext(seconds time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), seconds*time.Second)
}

func NewMongoClient(uri string, timeout time.Duration) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)
	timeoutContext, cancel := timeoutContext(timeout)
	defer cancel()

	client, err := mongo.Connect(timeoutContext, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(timeoutContext, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

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
	ctx, cancel := timeoutContext(m.timeout)
	defer cancel()

	entrySchema := entryFrom(entry)
	data, err := m.coll.InsertOne(ctx, entrySchema)

	if err != nil {
		return nil, err
	}

	if oid, ok := data.InsertedID.(primitive.ObjectID); ok {
		entrySchema.ID = oid
	}

	return entrySchema.toDomain(), nil
}

func (m *MongoEntryRepository) List() ([]domain.Entry, error) {
	ctx, cancel := timeoutContext(m.timeout)
	defer cancel()

	cursor, err := m.coll.Find(ctx, bson.D{})

	if err != nil {
		return nil, err
	}

	var list EntryListSchema

	for cursor.Next(ctx) {
		var entry EntrySchema
		err := cursor.Decode(&entry)
		if err != nil {
			return list.toDomain(), err
		}
		list = append(list, entry)
	}

	return list.toDomain(), nil
}
