FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /out/server /app/server

ENV HTTP_ADDR=:8080
ENV GREEN_API_BASE_URL=https://api.green-api.com

EXPOSE 8080

CMD ["/app/server"]
