FROM golang:1.24 AS builder

WORKDIR /src

# GOTOOLCHAIN=auto lets Go 1.24 automatically download the toolchain version
# declared in go.mod (e.g. go 1.25.0) instead of refusing to proceed.
ENV GOTOOLCHAIN=auto

# Copy go module files first for layer caching.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source and build a static binary.
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /e2e-example .

# ── Final stage ───────────────────────────────────────────────────────────────
FROM scratch

COPY --from=builder /e2e-example /e2e-example

EXPOSE 8080 8081

ENTRYPOINT ["/e2e-example"]
