# End-to-End Tracing Example

[![pkg.go.dev](https://pkg.go.dev/badge/go.opentelemetry.io/otel)](https://pkg.go.dev/go.opentelemetry.io/otel)

A self-contained example that shows the **Instana fork of opentelemetry-go**,
the **Instana fork of opentelemetry-go-contrib** HTTP instrumentation, and the
third-party **XSAM/otelsql** database instrumentation all composing into a
single distributed trace exported via **OTLP/gRPC**.

## What is demonstrated

| Component | Package | SDK type |
|---|---|---|
| Tracer setup, span SDK | `go.opentelemetry.io/otel/sdk` | **Core — Instana fork** |
| Semantic conventions | `go.opentelemetry.io/otel/semconv/v1.41.0` | **Core — Instana fork** |
| OTLP/gRPC exporter | `go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc` | **Core — Instana fork** |
| HTTP server + client middleware | `go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp` | **Contrib — Instana fork** |
| SQL instrumentation | `github.com/XSAM/otelsql` | **Third-party** |
| In-memory SQLite | `modernc.org/sqlite` | External (pure Go, no cgo) |

## Instana fork versions

All `go.opentelemetry.io` packages are redirected to the Instana public forks
via `replace` directives in `go.mod`:

| Upstream module | Instana fork | Version |
|---|---|---|
| `go.opentelemetry.io/otel` | `github.com/instana/opentelemetry-go` | `v0.1.0-otelv1.45.0` |
| `go.opentelemetry.io/otel/trace` | `github.com/instana/opentelemetry-go/trace` | `v0.1.0-otelv1.45.0` |
| `go.opentelemetry.io/otel/metric` | `github.com/instana/opentelemetry-go/metric` | `v0.1.0-otelv1.45.0` |
| `go.opentelemetry.io/otel/sdk` | `github.com/instana/opentelemetry-go/sdk` | `v0.1.0-otelv1.45.0` |
| `go.opentelemetry.io/otel/exporters/otlp/otlptrace` | `github.com/instana/opentelemetry-go/exporters/otlp/otlptrace` | `v0.1.0-otelv1.45.0` |
| `go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc` | `github.com/instana/opentelemetry-go/exporters/otlp/otlptrace/otlptracegrpc` | `v0.1.0-otelv1.45.0` |
| `go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp` | `github.com/instana/opentelemetry-go-contrib/instrumentation/net/http/otelhttp` | `v0.1.0-otelcontrib-ei-v0.70.0` |

Application imports remain unchanged (`go.opentelemetry.io/otel/…`); only
`go.mod` knows about the fork.

## Trace shape

For a single `GET /order/{id}` request the following spans are produced and sent to the collector:

```
GET /order/{id}                         ← otelhttp server span       (Instana contrib fork)
  ├── db: SELECT orders WHERE id = ?    ← otelsql child span         (third-party)
  └── GET http://localhost:8081/pay     ← otelhttp client span       (Instana contrib fork)
        └── payment-service             ← otelhttp server span       (Instana contrib fork)
```

The W3C `traceparent` header is injected by the order service's HTTP client and
extracted by the payment service middleware, linking all four spans into one
trace with a single root trace ID.

## Running with Docker Compose (recommended)

The easiest way to run the full stack — app, OTel Collector, and Jaeger — in
one command:

```bash
cd examples/e2e-example
docker compose up --build
```

This starts three services in dependency order:

| Service | Image | Purpose |
|---|---|---|
| `jaeger` | `jaegertracing/jaeger:latest` | Receives traces from the collector, UI on port `16686` |
| `collector` | `opentelemetry-collector-contrib:latest` | Receives OTLP/gRPC from the app, forwards to Jaeger, logs to stdout |
| `app` | Built from `Dockerfile` | The e2e-example binary (order service `:8080`, payment service `:8081`) |

The `Dockerfile` uses a multi-stage build: it compiles the Go source inside a
`golang:1.24` builder image (with `GOTOOLCHAIN=auto` to handle the `go 1.25.0`
toolchain directive) and copies the resulting static binary into a minimal
`scratch` final image.

Send a test request once all containers are up:

```bash
curl http://localhost:8080/order/42
```

Open Jaeger UI at [http://localhost:16686](http://localhost:16686) → select
service **order-service** → **Find Traces**.

To stop and remove all containers:

```bash
docker compose down
```

## Running locally (without Docker)

You still need a running OTel Collector. The quickest option:

```bash
docker run --rm \
  -p 4317:4317 \
  -p 4318:4318 \
  ghcr.io/open-telemetry/opentelemetry-collector-releases/opentelemetry-collector-contrib:latest
```

Then run the example:

```bash
cd examples/e2e-example
go run .
```

You should see:

```
payment service listening on :8081
order service  listening on :8080

Try:  curl http://localhost:8080/order/42
      curl http://localhost:8080/order/999  (not found)
```

## Sending test requests

```bash
# Successful order — produces a full 4-span trace
curl http://localhost:8080/order/42

# Another seeded order
curl http://localhost:8080/order/1

# Order not found — produces a 3-span trace (no payment call), order span has error status
curl http://localhost:8080/order/999
```

Expected JSON response for a successful order:

```json
{"amount":149,"customer":"Bob","item":"Gadget B","order_id":42,"status":"confirmed"}
```

## Viewing traces

### Jaeger UI (recommended)

When using Docker Compose, Jaeger is already running. Open:

[http://localhost:16686](http://localhost:16686) → select **order-service** → **Find Traces**

### Collector stdout

The collector logs received spans in detail to its own terminal (configured via
`otelcollector-config.yaml` with `verbosity: detailed`). Look for `"traceId"` —
all spans from one request share the same ID.

### Instana backend (SaaS)

Configure the OTel Collector to export to the Instana backend OTLP endpoint.
See [Sending OpenTelemetry data to the Instana backend](https://www.ibm.com/docs/en/instana-observability/saas?topic=instana-backend) for endpoint URLs, authentication (`x-instana-key`), and required resource attributes (`host.id` or `x-instana-host`).

Then view traces in Instana UI → **Traces** → search for **order-service**.

### Instana Agent

Run the Instana Host Agent with OTLP ingestion enabled (ports 4317/4318) and point the collector to it.
See [Sending OpenTelemetry data to the Instana agent](https://www.ibm.com/docs/en/iofgs?topic=instana-agent) for agent configuration and Kubernetes deployment.

Then view traces in Instana UI → **Traces**.

### Override the collector endpoint

```bash
OTEL_EXPORTER_OTLP_ENDPOINT=myhost:4317 go run .
```
