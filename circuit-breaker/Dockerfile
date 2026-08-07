FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

WORKDIR /app

RUN adduser -D -g '' appuser

COPY --from=builder /app/main .

USER appuser

EXPOSE 8080

CMD ["./main"]