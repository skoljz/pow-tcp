run-compose:
	docker-compose up -d --force-recreate

run:
	go run cmd/server/main.go

test:
	go clean --testcache
	go test ./...

deps:
	go mod download && go mod tidy

lint:
	GO111MODULE=on golangci-lint run