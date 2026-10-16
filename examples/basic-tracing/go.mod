module example.com/basic-tracing

go 1.25.0

require (
	go.opentelemetry.io/otel v1.43.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.33.0
	go.opentelemetry.io/otel/sdk v1.43.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	go.opentelemetry.io/otel/metric v1.43.0 // indirect
	go.opentelemetry.io/otel/trace v1.43.0 // indirect
	golang.org/x/sys v0.44.0 // indirect
)

// Replace standard OpenTelemetry with Instana fork (local development)
replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/exporters/stdout/stdoutlog => ../../exporters/stdout/stdoutlog
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric => ../../exporters/stdout/stdoutmetric
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace => ../../exporters/stdout/stdouttrace
	go.opentelemetry.io/otel/log => ../../log
	go.opentelemetry.io/otel/metric => ../../metric
	go.opentelemetry.io/otel/sdk => ../../sdk
	go.opentelemetry.io/otel/sdk/log => ../../sdk/log
	go.opentelemetry.io/otel/sdk/metric => ../../sdk/metric
	go.opentelemetry.io/otel/trace => ../../trace
)
