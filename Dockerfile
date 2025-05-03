FROM golang:1.23-alpine AS builder

RUN apk add --no-cache curl git

WORKDIR /app
COPY . .

RUN go build -o main main.go

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-arm64.tar.gz \
    | tar xvz -C /app

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY db/migration ./migration
COPY start.sh .
COPY wait-for.sh .

RUN chmod +x /app/start.sh /app/wait-for.sh /app/migrate

EXPOSE 8080

ENTRYPOINT ["/app/start.sh"]
CMD ["/app/main"]
