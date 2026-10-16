package main

import (
	"fmt"
	"net/http"
	"strconv"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// paymentHandler simulates a downstream payment microservice.
//
// In production this would be a separate binary. Here it runs on :8081 in the
// same process to keep the example self-contained.
//
// The otelhttp middleware that wraps this handler (wired in main.go) extracts
// the W3C traceparent header injected by the order service's HTTP client and
// continues the distributed trace, so the span created here appears as a child
// of the outbound client span from the order service.
func paymentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)

	q := r.URL.Query()

	orderID, err := strconv.ParseInt(q.Get("order_id"), 10, 64)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "bad order_id")
		http.Error(w, "bad order_id", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(q.Get("amount"), 64)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "bad amount")
		http.Error(w, "bad amount", http.StatusBadRequest)
		return
	}

	span.SetAttributes(
		attribute.Int64("payment.order_id", orderID),
		attribute.Float64("payment.amount", amount),
		attribute.String("payment.method", "credit_card"),
		attribute.Bool("payment.authorized", true),
	)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"authorized":true,"order_id":%d}`, orderID)
}
