FROM golang:1.18-bullseye AS builder

WORKDIR /go/src/github.com/upper-institute/flipbook/

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./flipbook github.com/upper-institute/flipbook/cmd

FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

WORKDIR /root/

COPY --from=builder /go/src/github.com/upper-institute/flipbook/flipbook /usr/local/bin/flipbook

EXPOSE 6336/tcp

ENTRYPOINT [ "flipbook" ]