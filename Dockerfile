# buidler
FROM golang:1.14-alpine

ENV GOOS=linux

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

## binaries
COPY pkg/geo/data/GeoLite2-Country.mmdb /data/GeoLite2-Country.mmdb