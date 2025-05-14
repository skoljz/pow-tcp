# build
FROM golang:1.22-alpine as builder
WORKDIR /src

COPY go.mod go.sum ./
RUN  go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./cmd/client

# run
FROM alpine:latest
WORKDIR /src
COPY --from=builder /src/app ./
CMD ["./app"]