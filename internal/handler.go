package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// allow "monkey patch" for unit tests
// this could go wrong with parallel tests
var now = time.Now

// default format will be RFC3339/ISO8601 which seems to be the preferred format
// for JSON according to most docs
func getTimestamp(format string) json.Marshaler {
	ts := now()
	ts = ts.UTC()
	if strings.ToLower(format) == "unix" {
		return UnixTimestamp(ts)
	}
	return ts
}

func DateTimeHandler(w http.ResponseWriter, r *http.Request) {
	ts := getTimestamp(r.URL.Query().Get("format"))

	result := struct {
		Timestamp interface{} `json:"timestamp"`
	}{
		ts,
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
