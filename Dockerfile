FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/server ./cmd/api

FROM alpine:3.20
COPY --from=builder /app/server /server
ENTRYPOINT ["/server"]
