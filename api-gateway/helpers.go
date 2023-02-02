package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (a *Config) readJson(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // one mb

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(data); err != nil {
		fmt.Println(err)
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("body must contain atlease one field")
	}
	return nil
}

func (a *Config) writeJson(w http.ResponseWriter, status int, data any, header ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if len(header) > 0 {
		for key, value := range header[0] {
			w.Header()[key] = value
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func (a *Config) errorJson(w http.ResponseWriter, err error, status ...int) error {
	statuscode := http.StatusBadRequest

	if len(status) > 0 {
		statuscode = status[0]
	}
	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()
	return a.writeJson(w, statuscode, payload)
}
