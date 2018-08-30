package main

import (
	"log"
	"net/http"

	"github.com/jbweber/project/v1/internal"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {

	inFlightGauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "in_flight_requests",
			Help: "A gauge of requests currently being served",
		},
	)

	durationHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "request_duration_seconds",
			Help: "A histogram of latencies for requests",
			//			Buckets: prometheus.DefBuckets,
			Buckets: []float64{0.005, 0.1, 0.5, 1},
		},
		[]string{"handler", "method"},
	)

	prometheus.MustRegister(inFlightGauge, durationHistogram)

	dateChain := promhttp.InstrumentHandlerInFlight(inFlightGauge,
		promhttp.InstrumentHandlerDuration(durationHistogram.MustCurryWith(prometheus.Labels{"handler": "DateTimeHandler"}), http.HandlerFunc(internal.DateTimeHandler)),
	)

	http.Handle("/", logRequest(http.HandlerFunc(internal.NotFoundHandler)))
	http.Handle("/datetime", logRequest(dateChain))
	http.Handle("/health", logRequest(http.HandlerFunc(internal.HealthHandler)))
	http.Handle("/info", logRequest(http.HandlerFunc(internal.InfoHandler)))
	http.Handle("/metrics", logRequest(promhttp.Handler()))

	log.Fatal(http.ListenAndServe(":8000", nil))
}
