# syntax=docker/dockerfile:1
FROM golang:1.25.1-alpine AS build

# Version can be set at build time
ARG VERSION=dev

WORKDIR /app

# Install git for modules
RUN apk add --no-cache git

# Copy module files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary with embedded version
ENV CGO_ENABLED=0
RUN go build -trimpath -ldflags="-s -w -X main.version=${VERSION}" -o /out/resolver ./cmd/resolver

# Final stage - using alpine for curl health checks
FROM alpine:3.20

RUN apk --no-cache add ca-certificates curl && \
    adduser -D app

USER app
WORKDIR /home/app

COPY --from=build /out/resolver /resolver

EXPOSE 8080

ENTRYPOINT ["/resolver", "--addr", ":8080"]