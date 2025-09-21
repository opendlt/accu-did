# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY . .

# Build the resolver
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o resolver ./cmd/resolver

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl
WORKDIR /root/

COPY --from=builder /app/resolver .

EXPOSE 8080

CMD ["./resolver", "--addr", ":8080", "--mode", "REAL"]