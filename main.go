package main

import (
	"log"
	"net/http"
	"nickel/app"
	"nickel/env"
	"nickel/http/handlers"
	"nickel/http/routers"
)

func run() error {
	env := env.NewEnv("./", "dev")
	app, err := app.NewApp(env.GetProp("DB_NAME"), env.GetProp("DB_URI"))

	if err != nil {
		return err
	}

	entryHandlers := handlers.NewEntryHandler(app.EntryService, app.Serializer)

	routers.ListenEntityRouters(app.Router, entryHandlers)
	return http.ListenAndServe(env.GetProp("PORT"), app.Router)
}

func main() {
	log.Fatal(run())
}
