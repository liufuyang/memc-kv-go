package metrics

import (
	"example.com/http_kv/cache"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

var (
	inFlightGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "in_flight_requests",
		Help: "A gauge of requests currently being served by the wrapped handler.",
	})

	counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_requests_total",
			Help: "A counter for requests to the wrapped handler.",
		},
		[]string{"code", "method"},
	)

	// duration is partitioned by the HTTP method and handler. It uses custom
	// buckets based on the expected request duration.
	duration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "A histogram of latencies for requests.",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"handler", "method", "cache"},
	)
	// responseSize has no labels, making it a zero-dimensional
	// ObserverVec.
	responseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "response_size_bytes",
			Help:    "A histogram of response sizes for requests.",
			Buckets: []float64{200, 500, 900, 1500},
		},
		[]string{},
	)
)

// InitMetrics can be called to start a promhttp for reporting metrics on an address
// For example
//   InitMetrics(":9000")
func InitMetrics(address string) {
	go func() {
		log.Printf("promhttp is listening on address: %v\n", address)
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(address, nil)
		if err != nil {
			panic(err)
		}
	}()
	// Register all the metrics in the standard registry.
	prometheus.MustRegister(inFlightGauge, counter, duration)
}

func CreateHttpHandleChain(f func(w http.ResponseWriter, req *http.Request),
	handlerLabel string,
	cacheLabel string,
	middleware func(f http.HandlerFunc) http.HandlerFunc) http.Handler {
	// Instrument the handlers with all the metrics, injecting the "handler"
	// label by currying.
	return promhttp.InstrumentHandlerInFlight(inFlightGauge,
		promhttp.InstrumentHandlerDuration(duration.MustCurryWith(prometheus.Labels{"handler": handlerLabel,
			"cache": cacheLabel}),
			promhttp.InstrumentHandlerCounter(counter,
				promhttp.InstrumentHandlerResponseSize(responseSize, middleware(f)),
			),
		),
	)
}

func RegisterCacheSizeGauge(c cache.Cache) {
	cacheSizeGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "cache_size_" + c.Name(),
		Help: "Size of the cache",
	})

	go func() {
		for {
			cacheSizeGauge.Set(float64(c.Size()))
			time.Sleep(5 * time.Second)
		}
	}()
}
