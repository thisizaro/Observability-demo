# --- build stage ---
# Full Go toolchain, only needed to produce the binary. This image
# itself is never shipped or run.
FROM golang:1.26-alpine AS builder
WORKDIR /app

# Copy just the dependency manifests first and download deps before
# copying the rest of the source. Docker caches each instruction as a
# layer; as long as go.mod/go.sum haven't changed, this layer is
# reused on the next build even if application code changed, so
# `go mod download` doesn't re-run on every single code edit.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

# --- final stage ---
# Alpine, not the builder image: the final container only needs to run
# a static binary, not compile one. Starting fresh here means the Go
# toolchain (hundreds of MB) never ends up in the image we actually
# deploy.
FROM alpine:3.20
COPY --from=builder /server /server

EXPOSE 8080
ENTRYPOINT ["/server"]
