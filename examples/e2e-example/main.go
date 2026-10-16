// Package main is an end-to-end tracing example that shows the Instana fork of
// opentelemetry-go (core SDK), the official opentelemetry-go-contrib HTTP
// instrumentation, and the third-party XSAM/otelsql database instrumentation
// all composing correctly within a single distributed trace.
//
// Architecture:
//
//	curl/browser
//	    │  GET /order/{id}
//	    ▼
//	[otelhttp middleware]  ← contrib SDK, server span
//	    │
//	    ├── db: SELECT orders  ← XSAM/otelsql, child span
//	    │
//	    └── GET /pay           ← otelhttp transport, child span
//	            └── [payment handler]  ← otelhttp middleware, linked span
//
// Traces are exported via OTLP/gRPC to the collector at localhost:4317.
// The collector endpoint can be overridden with the OTEL_EXPORTER_OTLP_ENDPOINT
// environment variable (standard OTel SDK env var).
//
// Run:
//
//	go run .
//	curl http://localhost:8080/order/42
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.41.0"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// collectorEndpoint is the default gRPC endpoint for the OTLP collector.
// Override at runtime with OTEL_EXPORTER_OTLP_ENDPOINT.
const collectorEndpoint = "localhost:4317"

func main() {
	ctx := context.Background()

	// ── Tracer setup ───────────────────────────────────────────────────────────

	// Resolve endpoint: honour the standard OTel env var if set.
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = collectorEndpoint
	}

	// Dial the collector. WithBlock() ensures we get a clear error at startup
	// rather than silently dropping spans until the connection comes up.
	conn, err := grpc.NewClient(endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("grpc.NewClient %s: %v", endpoint, err)
	}
	defer conn.Close()

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		log.Fatalf("otlptracegrpc.New: %v", err)
	}

	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName("order-service"),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		log.Fatalf("resource.New: %v", err)
	}

	// Instana fork of the core SDK — TracerProvider wired here.
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("TracerProvider.Shutdown: %v", err)
		}
	}()

	// Register globally so otelhttp and otelsql can reach the provider without
	// an explicit dependency on this main package.
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, // W3C traceparent/tracestate
		propagation.Baggage{},
	))

	// ── Database ───────────────────────────────────────────────────────────────

	db, err := initDB()
	if err != nil {
		log.Fatalf("initDB: %v", err)
	}
	defer db.Close()

	// ── HTTP servers ───────────────────────────────────────────────────────────

	// Payment service — separate port, simulates a downstream microservice.
	// otelhttp picks up the traceparent header and continues the trace.
	paymentMux := http.NewServeMux()
	paymentMux.Handle("/pay", otelhttp.NewHandler(
		http.HandlerFunc(paymentHandler),
		"payment-service",
		otelhttp.WithTracerProvider(tp),
		otelhttp.WithPropagators(otel.GetTextMapPropagator()),
	))
	paymentSrv := &http.Server{
		Addr:    ":8081",
		Handler: paymentMux,
	}

	// Order service — public port.
	// otelhttp creates the root server span for each incoming request.
	orderMux := http.NewServeMux()
	orderMux.Handle("/order/", otelhttp.NewHandler(
		newOrderHandler(db),
		"order-service",
		otelhttp.WithTracerProvider(tp),
		otelhttp.WithPropagators(otel.GetTextMapPropagator()),
	))
	orderSrv := &http.Server{
		Addr:    ":8080",
		Handler: orderMux,
	}

	// Start both servers.
	go func() {
		fmt.Println("payment service listening on :8081")
		if err := paymentSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("payment server: %v", err)
		}
	}()

	// Give the payment server a moment to bind before the order server starts
	// accepting traffic that would immediately call it.
	time.Sleep(50 * time.Millisecond)

	go func() {
		fmt.Println("order service  listening on :8080")
		fmt.Println()
		fmt.Println("Try:  curl http://localhost:8080/order/42")
		fmt.Println("      curl http://localhost:8080/order/999  (not found)")
		fmt.Println()
		if err := orderSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("order server: %v", err)
		}
	}()

	// ── Graceful shutdown ──────────────────────────────────────────────────────

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down…")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = orderSrv.Shutdown(ctx)
	_ = paymentSrv.Shutdown(ctx)
}
