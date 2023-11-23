package metric

import (
	"context"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"time"
)

type Meter struct {
	RequestCount    metric.Int64Counter
	RequestDuration metric.Float64Histogram
	Provider        *sdk.MeterProvider
}

func InitOTLPExporter(appName, otelExporterEndpoint string) (*Meter, error) {
	resources := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(appName),
	)

	exporter, err := otlpmetricgrpc.New(
		context.Background(),
		otlpmetricgrpc.WithTimeout(1*time.Second),
		otlpmetricgrpc.WithEndpoint(otelExporterEndpoint),
		otlpmetricgrpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	meterProvider := sdk.NewMeterProvider(
		sdk.WithResource(resources),
		sdk.WithReader(sdk.NewPeriodicReader(exporter)),
	)

	meter := meterProvider.Meter(
		appName,
		metric.WithInstrumentationVersion("v0.0.0"),
	)

	requestCount, err := meter.Int64Counter(
		"trace_app_otlp_request_count",
		metric.WithDescription("Incoming request count"),
		metric.WithUnit("request"),
	)
	if err != nil {
		return nil, err
	}

	requestDuration, err := meter.Float64Histogram(
		"trace_app_otlp_duration",
		metric.WithDescription("Incoming end to end duration"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return nil, err
	}

	return &Meter{
		RequestCount:    requestCount,
		RequestDuration: requestDuration,
		Provider:        meterProvider,
	}, nil
}
