# Build stage
FROM golang:1.24-alpine3.21 AS builder

RUN apk add --no-cache git ca-certificates make
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build 


# Runtime stage
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata curl
WORKDIR /app

COPY --from=builder /build/bin/swipe ./swipe

EXPOSE 8080
CMD ["./swipe"]