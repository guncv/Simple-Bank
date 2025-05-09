FROM golang:1.24-alpine AS builder

RUN apk add --no-cache curl git

WORKDIR /app
COPY . .

RUN go build -o main main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/main .
COPY app.env .
COPY db/migration ./migration
COPY wait-for.sh .

RUN chmod +x /app/wait-for.sh

EXPOSE 8080

CMD ["/app/wait-for.sh", "postgres12:5432", "--", "/app/main"]
