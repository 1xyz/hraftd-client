FROM golang:alpine AS builder

RUN apk update && apk add make git build-base && \
     rm -rf /var/cache/apk/*

ADD . /go/src/github.com/1xyz/hraftd-client
WORKDIR /go/src/github.com/1xyz/hraftd-client
RUN make release/linux

###

FROM alpine:latest AS hraftd-client

RUN apk update && apk add ca-certificates bash
WORKDIR /root/
COPY --from=builder /go/src/github.com/1xyz/hraftd-client/bin/linux/hraftc .