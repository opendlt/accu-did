# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY . .

# Build the registrar
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o registrar ./cmd/registrar

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl
WORKDIR /root/

COPY --from=builder /app/registrar .

EXPOSE 8081

CMD ["./registrar", "--addr", ":8081", "--mode", "REAL"]