package app

import (
	"nickel/core/services"
	"nickel/repository"
	"nickel/repository/config"
	"nickel/repository/mongo"
	"nickel/serializer"
	"nickel/serializer/json"

	"github.com/go-chi/chi/v5"
	m "go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	Router       *chi.Mux
	Mongo        *m.Client
	EntryRepo    repository.Entry
	Serializer   serializer.Serializer
	EntryService services.EntryService
}

func NewApp(dbName string, dbUri string) (*App, error) {
	r := chi.NewRouter()
	client, err := config.NewMongoClient(dbUri, 15)

	if err != nil {
		return nil, err
	}

	repoE := mongo.NewEntryRepository(client, dbName, 10)
	repoT := mongo.NewTagRepository(client, dbName, 10)
	serializer := json.NewJsonSerializer()
	serviceE := services.NewEntryService(repoE, repoT)

	return &App{
		Router:       r,
		Mongo:        client,
		EntryRepo:    repoE,
		Serializer:   serializer,
		EntryService: serviceE,
	}, nil
}
