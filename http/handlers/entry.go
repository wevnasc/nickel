package handlers

import (
	"net/http"
	"nickel/core/errors"
	"nickel/core/ports"
	"nickel/http/in"
	"nickel/http/out"
	"nickel/serializer/json"

	"github.com/go-chi/chi/v5"
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
	return ErrorHandler(h.serializer, func(w http.ResponseWriter, r *http.Request) error {
		entry := in.Entry{}
		err := json.DecodeBody(h.serializer, r.Body, &entry)

		if err != nil {
			return errors.Wrap(
				errors.Serialization,
				"not was possible decode entry invalid body",
				err,
			)
		}

		newEntry, err := h.service.Create(entry.Domain())
		if err != nil {
			return err
		}

		out := out.EntryFrom(newEntry)
		res, err := h.serializer.Encode(out)
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(res)
		return nil
	})
}

func (h *EntryHandlers) List() http.HandlerFunc {
	return ErrorHandler(h.serializer, func(w http.ResponseWriter, r *http.Request) error {
		entries, err := h.service.List()

		if err != nil {
			return err
		}

		out := out.EntriesListFrom(entries)
		res, err := h.serializer.Encode(out)

		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return nil
	})
}

func (h *EntryHandlers) Delete() http.HandlerFunc {
	return ErrorHandler(h.serializer, func(w http.ResponseWriter, r *http.Request) error {
		ID := chi.URLParam(r, "id")
		err := h.service.Delete(ID)

		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusNoContent)
		return nil
	})
}
