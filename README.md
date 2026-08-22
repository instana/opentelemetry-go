# Instana OpenTelemetry Go

> [!IMPORTANT]
> This project is currently in **beta** status and is susceptible to breaking changes. APIs, features, and functionality may change without notice. Use in production environments at your own risk.

## Overview

Instana OpenTelemetry Go is based on Open Source [OpenTelemetry Go](https://github.com/open-telemetry/opentelemetry-go). It provides an OpenTelemetry implementation which focuses on supporting Instana OpenTelemetry and IBM platforms (S390X Linux, PowerPC Linux, AIX) as well as other platforms (Linux x64/ARM64, macOS, and Windows).

It includes standard features which provide a collection of tools, APIs, and SDKs used to instrument, generate, collect, and export telemetry data (metrics, logs, and traces) for analysis in order to understand your software's performance and behavior.

In addition to standard features of OpenTelemetry Go, Instana OpenTelemetry Go also provides Instana OpenTelemetry Go SDK which supports exporter, propagation and serialization. It allows you to send OpenTelemetry trace data to Instana for monitoring and observability.

## Project Status

| Signal  | Status             |
|---------|--------------------|
| Traces  | Stable             |
| Metrics | Stable             |
| Logs    | Beta               |

Project versioning information and stability guarantees can be found in the [versioning documentation](VERSIONING.md).

### Compatibility

OpenTelemetry-Go ensures compatibility with the current supported versions of the [Go language](https://golang.org/doc/devel/release#policy):

> Each major Go release is supported until there are two newer major releases.
> For example, Go 1.5 was supported until the Go 1.7 release, and Go 1.6 was supported until the Go 1.8 release.

Currently, this project supports the following environments:

| OS       | Go Version | Architecture |
|----------|------------|--------------|
| Ubuntu   | 1.26       | amd64        |
| Ubuntu   | 1.25       | amd64        |
| Ubuntu   | 1.26       | 386          |
| Ubuntu   | 1.25       | 386          |
| Ubuntu   | 1.26       | arm64        |
| Ubuntu   | 1.25       | arm64        |
| macOS    | 1.26       | amd64        |
| macOS    | 1.25       | amd64        |
| macOS    | 1.26       | arm64        |
| macOS    | 1.25       | arm64        |
| Windows  | 1.26       | amd64        |
| Windows  | 1.25       | amd64        |
| Windows  | 1.26       | 386          |
| Windows  | 1.25       | 386          |

While this project should work for other systems, no compatibility guarantees are made for those systems currently.

## Getting Started

### Download and Build

The Instana OpenTelemetry Go is available in source as `tar.gz` or `zip` file which can be downloaded from releases.

Before building from source, make sure the following tools are installed:
- Go 1.25 or above
- Standard build tools for your platform (gcc/g++ on Linux/AIX, Xcode Command Line Tools on macOS, Visual Studio Build Tools on Windows)

To build everything, run:
```bash
make
```

### Using the Library

To use Instana OpenTelemetry Go in your project, add a replace directive in your `go.mod` file to point to the local path or repository:

```go
replace go.opentelemetry.io/otel => /path/to/instana-opentelemetry-go
replace go.opentelemetry.io/otel/sdk => /path/to/instana-opentelemetry-go/sdk
```

Then import the packages in your Go code as usual:
```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/sdk/trace"
)
```

OpenTelemetry's goal is to provide a single set of APIs to capture distributed traces and metrics from your application and send them to an observability platform. This project allows you to do just that for applications written in Go.

#### Instrumentation

To start capturing distributed traces and metric events from your application it first needs to be instrumented. The easiest way to do this is by using an instrumentation library for your code. Be sure to check out the officially supported instrumentation libraries.

If you need to extend the telemetry an instrumentation library provides or want to build your own instrumentation for your application directly you will need to use the [Go otel](https://pkg.go.dev/go.opentelemetry.io/otel) package.

#### Export

Now that your application is instrumented to collect telemetry, it needs an export pipeline to send that telemetry to an observability platform.

All officially supported exporters for the OpenTelemetry project are contained in the [exporters directory](./exporters).

| Exporter                              | Logs | Metrics | Traces |
|---------------------------------------|:----:|:-------:|:------:|
| [OTLP](./exporters/otlp/)             |  ✓   |    ✓    |   ✓    |
| [Prometheus](./exporters/prometheus/) |      |    ✓    |        |
| [stdout](./exporters/stdout/)         |  ✓   |    ✓    |   ✓    |
| [Zipkin](./exporters/zipkin/)         |      |         |   ✓    |

### Examples

Examples are provided in the [examples directory](./examples) to demonstrate how to use Instana OpenTelemetry Go:

- **[basic-tracing](./examples/basic-tracing)** - Shows how to create a tracer provider, set up service resource information, create parent and child spans, and export traces to stdout

Each example includes:
- Complete working code
- Instructions for local development using replace directives
- Instructions for using the published fork from GitHub
- Expected output and key concepts

To run an example locally:
```bash
cd examples/basic-tracing
go run main.go
```

For production use, update your `go.mod` to point to the published fork:
```go
replace go.opentelemetry.io/otel => github.com/instana/opentelemetry-go v1.33.0
replace go.opentelemetry.io/otel/sdk => github.com/instana/opentelemetry-go/sdk v1.43.0
```

## Contributing

See the [contributing documentation](CONTRIBUTING.md).
