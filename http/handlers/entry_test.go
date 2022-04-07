package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"nickel/app"
	"nickel/env"
	"nickel/http/in"
	"nickel/http/out"
	"nickel/repository/config"
	"nickel/repository/mongo"
	"nickel/serializer/json"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var testEnv *env.Env
var testApp *app.App
var entryHandlers *EntryHandlers

func TestMain(m *testing.M) {
	testEnv = env.NewEnv("../../", "test")
	testApp, _ = app.NewApp(testEnv.GetProp("DB_NAME"), testEnv.GetProp("DB_URI"))

	entryHandlers = NewEntryHandler(testApp.EntryService, testApp.Serializer)
	code := m.Run()
	os.Exit(code)
}

func cleanDatabase() {
	ctx, cancel := config.TimeoutContext(3)
	defer cancel()
	testApp.Mongo.Database(testEnv.GetProp("DB_NAME")).Collection("entries").Drop(ctx)
}

func createEntry() *mongo.Entry {
	coll := testApp.Mongo.Database(testEnv.GetProp("DB_NAME")).Collection("entries")
	entry := mongo.Entry{
		ID:          primitive.NewObjectID(),
		Description: "Ice cream",
		Amount:      4.5,
		Tags:        []mongo.EntryTag{{ID: primitive.NewObjectID(), Name: "Grocery"}},
		Type:        "Expense",
	}
	coll.InsertOne(context.Background(), entry)
	return &entry
}

func TestCreateEntryWithSuccess(t *testing.T) {
	payload := in.Entry{
		Description: "Ice cream",
		Amount:      4.5,
		Tags:        []string{"Grocery"},
		Type:        "Expense",
	}

	data, err := json.EncodeBody(testApp.Serializer, payload)

	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/", data)
	w := httptest.NewRecorder()

	createHandler := entryHandlers.Create()
	createHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	var entry out.Entry
	json.DecodeBody(testApp.Serializer, res.Body, &entry)

	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotNil(t, entry.ID)
	assert.Equal(t, payload.Description, entry.Description)
	assert.Equal(t, payload.Amount, entry.Amount)
	assert.Equal(t, payload.Type, entry.Type)
	assert.Equal(t, payload.Tags, entry.Tags)

	cleanDatabase()
}

func TestListEntriesWithSuccess(t *testing.T) {
	savedEntry := createEntry()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	listHandler := entryHandlers.List()
	listHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	var entries []out.Entry
	json.DecodeBody(testApp.Serializer, res.Body, &entries)

	assert.NotEmpty(t, entries)

	entry := entries[0]
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, savedEntry.ID.Hex(), entry.ID)
	assert.Equal(t, savedEntry.Description, entry.Description)
	assert.Equal(t, savedEntry.Amount, entry.Amount)
	assert.Equal(t, savedEntry.Type, entry.Type)

	for idx, tag := range savedEntry.Tags {
		assert.Equal(t, tag.Name, entry.Tags[idx])
	}

	cleanDatabase()
}

func TestEmptyEntriesWithSuccess(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	listHandler := entryHandlers.List()
	listHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	var entries []out.Entry
	json.DecodeBody(testApp.Serializer, res.Body, &entries)

	assert.Empty(t, entries)
	cleanDatabase()
}

func TestDeleteEntryWithSuccess(t *testing.T) {
	savedEntry := createEntry()
	ID := savedEntry.ID.Hex()
	uri := fmt.Sprintf("/%s", ID)
	req := httptest.NewRequest(http.MethodDelete, uri, nil)

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", ID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	w := httptest.NewRecorder()

	deleteHandler := entryHandlers.Delete()
	deleteHandler(w, req)

	res := w.Result()

	assert.Equal(t, http.StatusNoContent, res.StatusCode)
	cleanDatabase()
}

func TestDeleteEntryNotFound(t *testing.T) {
	ID := primitive.NewObjectID().Hex()
	uri := fmt.Sprintf("/%s", ID)
	req := httptest.NewRequest(http.MethodDelete, uri, nil)

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", ID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	w := httptest.NewRecorder()

	deleteHandler := entryHandlers.Delete()
	deleteHandler(w, req)

	res := w.Result()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	cleanDatabase()
}
