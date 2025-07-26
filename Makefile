swagger: 
	swag fmt
	swag init -g ./cmd/router.go
dev:
	docker-compose up --build 
run: 
	go run ./cmd/.
test:
	go test -v ./...
test-coverage:
	go test -v ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
unit-test:
	go test -v ./tests/unit/... 
