# base go image
FROM golang:1.18-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o frontEndApp ./cmd/web

RUN chmod +x /app/frontEndApp

# build tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/frontEndApp /app

CMD ["/app/frontEndApp"]