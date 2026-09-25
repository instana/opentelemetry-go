package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.41.0"
	"go.opentelemetry.io/otel/trace"
)

// paymentServiceURL is the address of the payment microservice running in the
// same process on a different port.
const paymentServiceURL = "http://localhost:8081/pay"

// orderHandler handles GET /order/{id}.
// It queries the database and then calls the payment service, demonstrating
// both database and outbound HTTP child spans under the same root server span.
type orderHandler struct {
	db         *sql.DB
	httpClient *http.Client
}

func newOrderHandler(db *sql.DB) *orderHandler {
	return &orderHandler{
		db: db,
		// otelhttp.NewTransport wraps the default transport so that every
		// outbound request carries a W3C traceparent header and produces a
		// child HTTP client span under the current context.
		httpClient: &http.Client{
			Transport: otelhttp.NewTransport(
				http.DefaultTransport,
				otelhttp.WithTracerProvider(otel.GetTracerProvider()),
				otelhttp.WithPropagators(otel.GetTextMapPropagator()),
			),
		},
	}
}

func (h *orderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract the order ID from the path: /order/{id}
	idStr := strings.TrimPrefix(r.URL.Path, "/order/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	// Add the order ID to the root server span created by otelhttp.
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.Int64("order.id", id))

	// ── Step 1: query the database ────────────────────────────────────────────
	// queryOrder passes ctx down to database/sql, which otelsql intercepts to
	// create a child "db: SELECT …" span automatically.
	order, err := queryOrder(ctx, h.db, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			span.SetAttributes(semconv.HTTPResponseStatusCode(http.StatusNotFound))
			http.Error(w, fmt.Sprintf("order %d not found", id), http.StatusNotFound)
			return
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, "db query failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	span.SetAttributes(
		attribute.String("order.customer", order.Customer),
		attribute.String("order.item", order.Item),
	)

	// ── Step 2: call the payment service ─────────────────────────────────────
	// The otelhttp transport injects traceparent into the outgoing request, so
	// the payment service continues the same distributed trace.
	payURL := fmt.Sprintf("%s?order_id=%d&amount=%.2f", paymentServiceURL, order.ID, order.Amount)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, payURL, nil)
	resp, err := h.httpClient.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "payment service call failed")
		http.Error(w, "payment service unreachable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		span.SetStatus(codes.Error, "payment declined")
		http.Error(w, "payment declined", http.StatusPaymentRequired)
		return
	}

	// ── Respond ───────────────────────────────────────────────────────────────
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"order_id": order.ID,
		"customer": order.Customer,
		"item":     order.Item,
		"amount":   order.Amount,
		"status":   "confirmed",
	})
}
