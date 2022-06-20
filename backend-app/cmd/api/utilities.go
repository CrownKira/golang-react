package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, wrap string) error {
	// wrap: wrap json with some kind of key
	wrapper := make(map[string]interface{}) // string to interface

	// {wrap: data}
	// convert wrapper to json byte slice
	wrapper[wrap] = data
	// MarshalIndent: will indent
	// Marshal: json in one line
	js, err := json.Marshal(wrapper)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	// write the status
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) errorJSON(w http.ResponseWriter, err error) {
	type jsonError struct {
		Message string `json:"message"`
	}

	theError := jsonError{
		// Error() gives the error message in string
		Message: err.Error(),
	}

	// write error json to writer; in json with the key error
	app.writeJSON(w, http.StatusBadRequest, theError, "error")

}
