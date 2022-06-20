package main

import (
	"encoding/json"
	"net/http"
)

// write to http channel the response
// use receiver application
// application has the statushandler, config, etc
// easy way to share things in clean and logical fashion
func (app *application) statusHandler(w http.ResponseWriter, r *http.Request) {
	// print to response writer
	// qn: print 2nd arg to w?
	// fmt.Fprint(w, "status")
	currentStatus := AppStatus{
		Status:      "Available",
		Environment: app.config.env,
		Version:     version,
	}

	// convert struct to json in byte slice, prefix, how much to indent
	js, err := json.MarshalIndent(currentStatus, "", "\t")
	if err != nil {
		app.logger.Println(err)
	}

	// set the header of writer: set to application/json
	w.Header().Set("Content-Type", "application/json")
	// write header status ok
	// send status
	w.WriteHeader(http.StatusOK)
	// write json to http channel
	w.Write(js)
}
