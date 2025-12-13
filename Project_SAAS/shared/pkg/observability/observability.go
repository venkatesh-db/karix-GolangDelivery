package observability

import (
	"context"
	"os"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.uber.org/zap"
)

// ShutdownFunc releases telemetry exporters.
type ShutdownFunc func(ctx context.Context) error

// Setup initializes OpenTelemetry tracing and propagation.
func Setup(ctx context.Context, service string, log *zap.Logger) (ShutdownFunc, error) {
	if strings.EqualFold(os.Getenv("OBSERVABILITY_DISABLED"), "true") {
		log.Info("observability disabled via env flag")
		return func(context.Context) error { return nil }, nil
	}

	exporter, exporterName, err := buildExporter(ctx)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithOS(),
		resource.WithProcess(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceName(service),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter,
			sdktrace.WithBatchTimeout(5*time.Second),
		),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	log.Info("observability initialized", zap.String("exporter", exporterName))

	return tp.Shutdown, nil
}

func buildExporter(ctx context.Context) (sdktrace.SpanExporter, string, error) {
	endpoint := strings.TrimSpace(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if endpoint != "" {
		normalized, insecure := normalizeEndpoint(endpoint)
		path := strings.TrimPrefix(strings.TrimSpace(os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT")), "/")
		if path == "" {
			path = "v1/traces"
		}
		opts := []otlptracehttp.Option{
			otlptracehttp.WithEndpoint(normalized),
			otlptracehttp.WithURLPath(path),
		}
		if insecure {
			opts = append(opts, otlptracehttp.WithInsecure())
		}
		client := otlptracehttp.NewClient(opts...)
		exporter, err := otlptrace.New(ctx, client)
		if err != nil {
			return nil, "otlphttp", err
		}
		return exporter, "otlphttp", nil
	}

	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, "stdout", err
	}
	return exporter, "stdout", nil
}

func normalizeEndpoint(endpoint string) (string, bool) {
	e := strings.TrimSpace(endpoint)
	insecure := true
	switch {
	case strings.HasPrefix(e, "https://"):
		e = strings.TrimPrefix(e, "https://")
		insecure = false
	case strings.HasPrefix(e, "http://"):
		e = strings.TrimPrefix(e, "http://")
	}
	e = strings.TrimSuffix(e, "/")
	return e, insecure
}
