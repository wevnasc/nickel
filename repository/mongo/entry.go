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

type EntryTag struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}

type EntryTagList []EntryTag

type Entry struct {
	ID          primitive.ObjectID `bson:"_id"`
	Description string             `bson:"description"`
	Amount      float64            `bson:"amount"`
	Tags        EntryTagList       `bson:"tags"`
	Type        string             `bson:"type"`
}

type EntryList []Entry

func entryFrom(entry *domain.Entry) *Entry {
	ID, _ := primitive.ObjectIDFromHex(entry.ID)

	if ID == primitive.NilObjectID {
		ID = primitive.NewObjectID()
	}

	return &Entry{
		ID:          ID,
		Description: entry.Description,
		Amount:      entry.Amount,
		Tags:        entryTagsFrom(entry.Tags),
		Type:        string(entry.Type),
	}

}

func (s *Entry) toDomain() *domain.Entry {
	return &domain.Entry{
		ID:          s.ID.Hex(),
		Description: s.Description,
		Amount:      s.Amount,
		Tags:        s.Tags.toDomain(),
		Type:        domain.Type(s.Type),
	}
}

func (s EntryList) toDomain() []domain.Entry {
	list := make([]domain.Entry, len(s))
	for idx, e := range s {
		entry := e.toDomain()
		list[idx] = *entry
	}

	return list
}

func entryTagFrom(tag *domain.Tag) *EntryTag {
	ID, _ := primitive.ObjectIDFromHex(tag.ID)

	if ID == primitive.NilObjectID {
		ID = primitive.NewObjectID()
	}

	return &EntryTag{
		ID:   ID,
		Name: tag.Name,
	}
}

func (s *EntryTag) toDomain() *domain.Tag {
	return &domain.Tag{
		ID:   s.ID.Hex(),
		Name: s.Name,
	}
}

func entryTagsFrom(tags []domain.Tag) EntryTagList {
	list := make(EntryTagList, len(tags))

	for idx, t := range tags {

		list[idx] = *entryTagFrom(&t)
	}

	return list
}

func (s EntryTagList) toDomain() []domain.Tag {
	list := make([]domain.Tag, len(s))
	for idx, t := range s {
		tag := t.toDomain()
		list[idx] = *tag
	}
	return list
}

type EntryRepository struct {
	timeout time.Duration
	client  *mongo.Client
	coll    *mongo.Collection
}

func NewEntryRepository(client *mongo.Client, database string, timeout time.Duration) *EntryRepository {
	return &EntryRepository{
		client:  client,
		timeout: timeout,
		coll:    client.Database(database).Collection("entries"),
	}
}

func (m *EntryRepository) Create(entry *domain.Entry) (*domain.Entry, error) {
	ctx, cancel := config.TimeoutContext(m.timeout)
	defer cancel()

	schema := entryFrom(entry)
	data, err := m.coll.InsertOne(ctx, entry)

	if err != nil {
		return nil, errors.Wrap(errors.InsertData, "not was possible to create the entry", err)
	}

	if oid, ok := data.InsertedID.(primitive.ObjectID); ok {
		schema.ID = oid
	}

	return schema.toDomain(), nil
}

func (m *EntryRepository) List() ([]domain.Entry, error) {
	ctx, cancel := config.TimeoutContext(m.timeout)
	defer cancel()

	cursor, err := m.coll.Find(ctx, bson.D{})

	if err != nil {
		return nil, errors.Wrap(errors.FindData, "not was possible to query the entries", err)
	}

	var list EntryList

	for cursor.Next(ctx) {
		var entry Entry
		err := cursor.Decode(&entry)
		if err != nil {
			return list.toDomain(), errors.Wrap(errors.Serialization, "not was to decode the entry from database", err)
		}
		list = append(list, entry)
	}

	return list.toDomain(), nil
}

func (m *EntryRepository) Get(id string) (*domain.Entry, error) {
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

	entry := &Entry{}
	err = res.Decode(entry)

	if err != nil {
		return nil, errors.Wrap(errors.Serialization, "not was possible decode db entry", err)
	}

	return entry.toDomain(), nil
}

func (m *EntryRepository) Delete(id string) error {
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
