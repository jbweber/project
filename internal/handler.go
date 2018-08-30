package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// allow "monkey patch" for unit tests
// this could go wrong with parallel tests
var now = time.Now

func DateTimeHandler(w http.ResponseWriter, r *http.Request) {
	ts := now()
	uts := UnixTimestamp(ts)

	result := struct {
		Timestamp interface{} `json:"timestamp"`
	}{
		uts,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&result)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

	// realistically our health check would be split into health and liveness
	// and would actually check something. our app doesn't do anything really
	// so this should be okay for now

	// implicitly should return OK with empty body
}

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	info := struct {
		GitCommit string `json:"gitCommit"`
		Version   string `json:"version"`
	}{GitCommit, Version}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&info)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, fmt.Sprintf("path %v not found", r.URL.Path), http.StatusNotFound)
}
