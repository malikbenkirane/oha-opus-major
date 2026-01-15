package otel

import (
	"fmt"
	"time"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

type Config struct {
	batcherTimeout    time.Duration
	traceSpanExporter trace.SpanExporter

	metricExporterInterval time.Duration
	metricExporter         metric.Exporter

	logExporter log.Exporter
}

type Option func(Config) Config

func defaultConfig() (Config, error) {
	traceExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return Config{}, fmt.Errorf("trace exporter: %w", err)
	}
	metricExporter, err := stdoutmetric.New()
	if err != nil {
		return Config{}, fmt.Errorf("metric exporter: %w", err)
	}
	logExporter, err := stdoutlog.New()
	if err != nil {
		return Config{}, fmt.Errorf("log exporter: %w", err)
	}
	return Config{
		batcherTimeout:    time.Second,
		traceSpanExporter: traceExporter,
		metricExporter:    metricExporter,
		logExporter:       logExporter,
	}, nil
}

func WithTraceSpanExporter(exporter trace.SpanExporter) Option {
	return func(c Config) Config {
		c.traceSpanExporter = exporter
		return c
	}
}

func WithBatcherTimeout(timeout time.Duration) Option {
	return func(c Config) Config {
		c.batcherTimeout = timeout
		return c
	}
}

func WithMetricExporter(exporter metric.Exporter) Option {
	return func(c Config) Config {
		c.metricExporter = exporter
		return c
	}
}

func WithMetricExporterInterval(interval time.Duration) Option {
	return func(c Config) Config {
		c.metricExporterInterval = interval
		return c
	}
}

type otel struct {
	config Config
}

func newOtel(opts ...Option) (*otel, error) {
	config, err := defaultConfig()
	if err != nil {
		return nil, err
	}
	for _, opt := range opts {
		config = opt(config)
	}
	return &otel{
		config: config,
	}, nil
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func (ot otel) newTraceProvider() *trace.TracerProvider {
	return trace.NewTracerProvider(
		trace.WithBatcher(
			ot.config.traceSpanExporter,
			trace.WithBatchTimeout(ot.config.batcherTimeout)))
}

func (ot otel) newMeterProvider() *metric.MeterProvider {
	return metric.NewMeterProvider(
		metric.WithReader(
			metric.NewPeriodicReader(
				ot.config.metricExporter,
				metric.WithInterval(ot.config.metricExporterInterval))))
}

func (o otel) newLoggerProvider() *log.LoggerProvider {
	return log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(o.config.logExporter)))
}
