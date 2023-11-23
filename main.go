package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"log"
	"net/http"
	"source.golabs.io/engineering-platforms/lens/trace-app-golang/config"
	"source.golabs.io/engineering-platforms/lens/trace-app-golang/environment"
	"source.golabs.io/engineering-platforms/lens/trace-app-golang/handler"
	localmetric "source.golabs.io/engineering-platforms/lens/trace-app-golang/metric"
	"strconv"
	"time"
)

func main() {
	time.Sleep(5 * time.Second)
	cfg := config.Load()
	env := environment.Init(cfg)
	defer shutdown(env, context.Background())

	r := mux.NewRouter()

	r.HandleFunc("/v1/books", handler.CreateBook(env)).Methods("POST")
	r.HandleFunc("/v1/books/{title}", handler.ReadBook(env)).Methods("GET")
	r.Use(otelmux.Middleware(cfg.AppName, otelmux.WithPropagators(propagation.NewCompositeTextMapPropagator(
		propagation.Baggage{},
		propagation.TraceContext{},
	))))
	r.Use(instrumentationMiddleware(env, cfg))

	if !cfg.OTLPMetricExport {
		go startMonitoringServer()
	}

	log.Printf("Listening on port 8080...")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Printf("Shutting down server")
	}
}

type CustomResponseWriter struct {
	StatusCode int
	Body       []byte
	http.ResponseWriter
}

// WriteHeader - a custom WriteHeader
func (crw *CustomResponseWriter) WriteHeader(code int) {
	crw.StatusCode = code
	crw.ResponseWriter.WriteHeader(code)
}

// Write - a custom Write
func (crw *CustomResponseWriter) Write(body []byte) (int, error) {
	crw.Body = body
	return crw.ResponseWriter.Write(body)
}

func instrumentationMiddleware(env environment.Environment, cfg config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.Background()
			requestStartTime := time.Now()
			crw := &CustomResponseWriter{
				ResponseWriter: w,
			}

			next.ServeHTTP(crw, r)

			elapsedTime := float64(time.Since(requestStartTime)) / float64(time.Millisecond)

			if cfg.OTLPMetricExport {
				attrs := semconv.HTTPServerMetricAttributesFromHTTPRequest(cfg.AppName, r)
				attrs = append(attrs, attribute.String("http_route", r.URL.Path))
				attrs = append(attrs, attribute.Int("http_status", crw.StatusCode))
				env.Meter.RequestCount.Add(ctx, 1, metric.WithAttributes(attrs...))
				env.Meter.RequestDuration.Record(ctx, elapsedTime, metric.WithAttributes(attrs...))
			} else {
				localmetric.RequestCount.WithLabelValues(strconv.Itoa(crw.StatusCode), r.Method, r.URL.Path).Inc()
				localmetric.RequestDuration.WithLabelValues(strconv.Itoa(crw.StatusCode), r.Method, r.URL.Path).Observe(elapsedTime)
			}
		})
	}
}

func shutdown(env environment.Environment, ctx context.Context) {
	env.TraceExporter.Shutdown(ctx)
	env.Meter.Provider.Shutdown(ctx)
}

func startMonitoringServer() {
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatalf("Error while starting monitoring server on port 9090 - %v", err)
	}
}
