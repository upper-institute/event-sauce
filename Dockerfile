FROM golang:1.18-bullseye AS builder

WORKDIR /go/src/github.com/upper-institute/event-sauce/

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./paprika github.com/upper-institute/event-sauce/cmd/paprika

FROM debian:bullseye-slim

WORKDIR /root/

COPY --from=builder /go/src/github.com/upper-institute/event-sauce/paprika /usr/local/bin/paprika

EXPOSE 6336/tcp