FROM golang:1.18-bullseye AS builder

WORKDIR /go/src/github.com/vflopes/chronon/

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./cmd/chronon github.com/vflopes/chronon/cmd

FROM debian:bullseye-slim

WORKDIR /root/

COPY --from=builder /go/src/github.com/vflopes/chronon/cmd/chronon/cmd ./chronon

EXPOSE 57755/tcp

ENTRYPOINT ["./eventsauce"]