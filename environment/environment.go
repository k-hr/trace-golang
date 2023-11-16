package environment

import (
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"log"
	"source.golabs.io/engineering-platforms/lens/trace-app-golang/config"
	"source.golabs.io/engineering-platforms/lens/trace-app-golang/db"
	"source.golabs.io/engineering-platforms/lens/trace-app-golang/tracing"
)

type Environment struct {
	DBProxy       *db.Proxy
	TraceExporter *otlptrace.Exporter
}

func Init(cfg config.Config) Environment {
	databaseProxy := db.New(cfg)
	traceExporter, err := tracing.InitTraceExporter(cfg.AppName, cfg.OTLPExporterEndpoint)
	if err != nil {
		log.Printf("Error while initialising trace exporter - %v\n", err)
	}

	return Environment{DBProxy: databaseProxy, TraceExporter: traceExporter}
}
