# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o ticket-system ./cmd/main.go

# Run stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/ticket-system .

ENV PORT=8080
ENV DATABASE_PATH=/root/ticket_system.db
ENV JWT_SECRET=change-me-in-production

EXPOSE 8080

CMD ["./ticket-system"]
