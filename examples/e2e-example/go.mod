module example.com/e2e-example

go 1.25.0

require (
	github.com/XSAM/otelsql v0.41.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.69.0
	go.opentelemetry.io/otel v1.45.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.44.0
	go.opentelemetry.io/otel/sdk v1.45.0
	go.opentelemetry.io/otel/trace v1.45.0
	google.golang.org/grpc v1.83.2
	modernc.org/sqlite v1.56.0
)

require (
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/felixge/httpsnoop v1.1.0 // indirect
	github.com/go-logr/logr v1.4.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.29.0 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/ncruces/go-strftime v1.0.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.45.0 // indirect
	go.opentelemetry.io/otel/metric v1.45.0 // indirect
	go.opentelemetry.io/proto/otlp v1.11.0 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260803160001-6ac0973c030d // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260803160001-6ac0973c030d // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	modernc.org/libc v1.74.4 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.11.0 // indirect
)

// Replace core OpenTelemetry packages with the Instana fork.
replace (
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp => github.com/instana/opentelemetry-go-contrib/instrumentation/net/http/otelhttp v0.1.0-otelcontrib-ei-v0.70.0
	go.opentelemetry.io/otel => github.com/instana/opentelemetry-go v0.1.0-otelv1.45.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace => github.com/instana/opentelemetry-go/exporters/otlp/otlptrace v0.1.0-otelv1.45.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc => github.com/instana/opentelemetry-go/exporters/otlp/otlptrace/otlptracegrpc v0.1.0-otelv1.45.0
	go.opentelemetry.io/otel/metric => github.com/instana/opentelemetry-go/metric v0.1.0-otelv1.45.0
	go.opentelemetry.io/otel/sdk => github.com/instana/opentelemetry-go/sdk v0.1.0-otelv1.45.0
	go.opentelemetry.io/otel/trace => github.com/instana/opentelemetry-go/trace v0.1.0-otelv1.45.0
)
