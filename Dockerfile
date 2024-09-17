FROM golang:1.22.1-alpine AS base

RUN apk update && apk add --no-cache ca-certificates

FROM base AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

FROM alpine:3.13 AS final
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .

EXPOSE 8080

CMD ["./main"]

