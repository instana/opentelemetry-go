package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
)

func main() {
	// Create a stdout exporter for demonstration
	exporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}

	// Create a resource with service information
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName("example-service"),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}

	// Create a tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// Set the global tracer provider
	otel.SetTracerProvider(tp)

	// Get a tracer
	tracer := tp.Tracer("example-tracer")

	// Create some spans
	ctx := context.Background()
	
	// Parent span
	ctx, parentSpan := tracer.Start(ctx, "parent-operation")
	parentSpan.SetAttributes(
		attribute.String("operation.type", "example"),
		attribute.Int("operation.count", 1),
	)
	
	// Simulate some work
	time.Sleep(100 * time.Millisecond)
	
	// Child span
	_, childSpan := tracer.Start(ctx, "child-operation")
	childSpan.SetAttributes(
		attribute.String("child.type", "processing"),
		attribute.Bool("child.success", true),
	)
	
	// Simulate child work
	time.Sleep(50 * time.Millisecond)
	childSpan.End()
	
	// Another child span
	_, child2Span := tracer.Start(ctx, "child-operation-2")
	child2Span.SetAttributes(
		attribute.String("child.type", "validation"),
		attribute.Int("items.validated", 42),
	)
	time.Sleep(30 * time.Millisecond)
	child2Span.End()
	
	parentSpan.End()

	fmt.Println("Example completed successfully!")
	fmt.Println("Traces exported to stdout above")
	fmt.Println("Using Instana fork via replace directives")
}