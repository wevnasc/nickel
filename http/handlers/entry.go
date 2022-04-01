package handlers

import (
	"log"
	"net/http"
	"nickel/core/ports"
	"nickel/http/in"
	"nickel/http/out"
	"nickel/serializer/json"
)

type EntryHandlers struct {
	service    ports.EntryServicePort
	serializer ports.SerializerPort
}

func NewEntryHandler(service ports.EntryServicePort, serializer ports.SerializerPort) *EntryHandlers {
	return &EntryHandlers{
		service:    service,
		serializer: serializer,
	}
}

func (h *EntryHandlers) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entry := in.Entry{}
		err := json.DecodeBody(h.serializer, r.Body, &entry)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newEntry, err := h.service.Create(entry.Domain())

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		out := out.EntryFrom(newEntry)
		res, err := h.serializer.Encode(out)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(res)

	}
}

func (h *EntryHandlers) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entries, err := h.service.List()

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		out := out.EntriesListFrom(entries)
		res, err := h.serializer.Encode(out)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
}
