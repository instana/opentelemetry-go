# Basic Tracing Example

This example demonstrates how to use the Instana OpenTelemetry Go fork for basic distributed tracing.

## Overview

This example shows:
- Creating a tracer provider with stdout exporter
- Setting up service resource information
- Creating parent and child spans
- Adding attributes to spans
- Proper span lifecycle management

## Running the Example

### Local Development (using local fork)
```bash
cd examples/basic-tracing
go run main.go
```

### Using Published Fork
Update `go.mod` to point to the published fork on GitHub:

```go
// Replace local paths with GitHub repository
replace (
    go.opentelemetry.io/otel => github.com/instana/opentelemetry-go v1.33.0
    go.opentelemetry.io/otel/exporters/stdout/stdouttrace => github.com/instana/opentelemetry-go/exporters/stdout/stdouttrace v1.33.0
    go.opentelemetry.io/otel/log => github.com/instana/opentelemetry-go/log v1.43.0
    go.opentelemetry.io/otel/metric => github.com/instana/opentelemetry-go/metric v1.43.0
    go.opentelemetry.io/otel/sdk => github.com/instana/opentelemetry-go/sdk v1.43.0
    go.opentelemetry.io/otel/sdk/log => github.com/instana/opentelemetry-go/sdk/log v1.43.0
    go.opentelemetry.io/otel/sdk/metric => github.com/instana/opentelemetry-go/sdk/metric v1.43.0
    go.opentelemetry.io/otel/trace => github.com/instana/opentelemetry-go/trace v1.43.0
)
```

## Expected Output

The program will:
1. Create a parent span named "parent-operation"
2. Create two child spans: "child-operation" and "child-operation-2"
3. Export all spans to stdout in JSON format
4. Display a success message

You should see JSON output showing the trace hierarchy with:
- Service name: "example-service"
- Service version: "1.0.0"
- Parent span with 2 child spans
- Various attributes on each span

## Key Concepts

### Tracer Provider Setup
```go
tp := sdktrace.NewTracerProvider(
    sdktrace.WithBatcher(exporter),
    sdktrace.WithResource(res),
)
otel.SetTracerProvider(tp)
```

### Creating Spans
```go
ctx, span := tracer.Start(ctx, "operation-name")
defer span.End()
```

### Adding Attributes
```go
span.SetAttributes(
    attribute.String("key", "value"),
    attribute.Int("count", 42),
    attribute.Bool("success", true),
)
```

## Fork Information

This example uses the Instana fork of OpenTelemetry Go:
- **Repository:** `github.com/instana/opentelemetry-go`
- **Module Path:** `go.opentelemetry.io/otel` (unchanged from upstream)
- **Usage:** Standard OpenTelemetry imports with replace directives

### Why Replace Directives?

The fork maintains the original module path `go.opentelemetry.io/otel` to ensure:
- No code changes needed in your application
- Seamless switching between upstream and fork
- Standard Go module fork pattern
- Internal packages remain properly isolated

### Client Usage Pattern

In your application's `go.mod`:

```go
module myapp

require (
    go.opentelemetry.io/otel v1.33.0
    go.opentelemetry.io/otel/sdk v1.33.0
)

// Point to Instana's fork
replace go.opentelemetry.io/otel => github.com/instana/opentelemetry-go v1.33.0
replace go.opentelemetry.io/otel/sdk => github.com/instana/opentelemetry-go/sdk v1.33.0
```

Your Go code remains unchanged:
```go
import "go.opentelemetry.io/otel"
```