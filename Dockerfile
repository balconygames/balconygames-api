# buidler
FROM golang:1.14-alpine as builder

ENV GOOS=linux

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -v -o main ./cmd/...

# release image

FROM gcr.io/distroless/static
COPY --from=builder /app/main /

## migrations
COPY modules/primary/migrations modules/primary/migrations
COPY modules/auth/migrations modules/auth/migrations
COPY modules/leaderboard/migrations modules/leaderboard/migrations

## binaries
COPY pkg/geo/data/GeoLite2-Country.mmdb /data/GeoLite2-Country.mmdb

CMD ["/main"]
