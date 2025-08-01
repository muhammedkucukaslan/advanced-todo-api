swagger: 
	swag fmt
	swag init -g ./cmd/router.go
	go run  ./cmd/.
run: 
	go run ./cmd/.
test:
	go test -v ./...
test-coverage:
	go test -v ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

