FROM golang:1.16.1-alpine AS builder

# INSTALL PACKAGES
RUN apk update && apk add --no-cache musl-dev gcc build-base libc-dev curl git ca-certificates &&  \
    update-ca-certificates

ENV GOPATH="$HOME/go"

WORKDIR $GOPATH/src

COPY . $GOPATH/src

RUN GOOS=linux go build -ldflags '-linkmode=external' -o /go/bin/websocket cmd/api/main.go
RUN GOOS=linux go build -ldflags '-linkmode=external' -o /go/bin/consumer cmd/consumer/main.go

FROM alpine:3.13

# UPDATE APK CACHE AND INSTALL PACKAGES | CONFIGURE TIMEZONE | INSTALL AWS DEPS
RUN apk update && apk upgrade && apk add --no-cache \
    tzdata ca-certificates && \
    cp /usr/share/zoneinfo/America/Sao_Paulo /etc/localtime; echo "America/Sao_Paulo" > /etc/timezone

# Copy assets
WORKDIR /app/assets
COPY --from=builder /go/src/assets/ .

# Copy static executables
WORKDIR /app
COPY --from=builder /go/bin/websocket .
COPY --from=builder /go/bin/consumer .