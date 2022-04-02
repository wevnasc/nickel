package handlers

import (
	"log"
	"net/http"
	e "nickel/core/errors"
	"nickel/core/ports"
)

type ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request) error

type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var ErrorStatusCode = map[e.ErrorType]int{
	e.NotFound:      http.StatusNotFound,
	e.InsertData:    http.StatusBadRequest,
	e.FindData:      http.StatusBadRequest,
	e.Serialization: http.StatusBadRequest,
}

func ErrorHandler(serializer ports.SerializerPort, handler ErrorHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)

		if err == nil {
			return
		}

		if appErr, ok := err.(*e.AppError); ok {
			log.Println(appErr)

			body := &HttpError{
				Code:    ErrorStatusCode[appErr.Type],
				Message: appErr.Message,
			}

			res, _ := serializer.Encode(body)
			w.WriteHeader(ErrorStatusCode[appErr.Type])
			w.Write(res)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
	}
}
