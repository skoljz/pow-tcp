# build
FROM golang:1.22-alpine as builder
WORKDIR /src

COPY go.mod go.sum ./
RUN  go mod download

COPY .. .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./cmd/server

# run
FROM alpine:latest
WORKDIR /src
COPY --from=builder /src/app         ./app
COPY --from=builder /src/quotes.txt ./quotes.txt
ENV STORAGE_FILE=/src/quotes.txt
CMD ["./app"]