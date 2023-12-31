package util

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

func Send(w *http.ResponseWriter, data any) {
    json.NewEncoder(*w).Encode(data)
}

func Receive[T any](body io.ReadCloser, cl *T) error {
	decoder := json.NewDecoder(body)

	if err := decoder.Decode(cl); err != nil {
		return errors.New("invalid parameters")
	}

	return nil
}

func Error(w *http.ResponseWriter, status int, message string) {
    (*w).WriteHeader(status)
    log.Println(message)
}
