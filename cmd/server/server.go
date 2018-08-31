package main

import (
	"log"
	"net/http"

	"github.com/jbweber/project/internal"
)

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {

	http.Handle("/", logRequest(http.HandlerFunc(internal.NotFoundHandler)))
	http.Handle("/datetime", logRequest(http.HandlerFunc(internal.DateTimeHandler)))
	http.Handle("/health", logRequest(http.HandlerFunc(internal.HealthHandler)))
	http.Handle("/info", logRequest(http.HandlerFunc(internal.InfoHandler)))

	log.Fatal(http.ListenAndServe(":8000", nil))
}
