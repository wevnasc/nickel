package main

import (
	"log"
	"net/http"
	"nickel/core/services"
	"nickel/env"
	"nickel/http/handlers"
	"nickel/http/routers"
	"nickel/repositories/config"
	"nickel/repositories/mongo/entry"
	"nickel/repositories/mongo/tag"
	"nickel/serializer/json"

	"github.com/go-chi/chi/v5"
)

func run() error {
	r := chi.NewRouter()
	env := env.NewEnv("dev")
	client, err := config.NewMongoClient(env.GetProp("DB_URI"), 15)

	if err != nil {
		log.Println("error to connect with mongo db")
		return err
	}

	repoE := entry.NewMongoEntryRepository(client, env.GetProp("DB_NAME"), 10)
	repoT := tag.NewMongoTagRepository(client, env.GetProp("DB_NAME"), 10)
	serializer := json.NewJsonSerializer()
	service := services.NewEntryService(repoE, repoT)
	handlers := handlers.NewEntryHandler(service, serializer)

	routers.ListenEntityRouters(r, handlers)
	return http.ListenAndServe(env.GetProp("PORT"), r)
}

func main() {
	log.Fatal(run())
}
