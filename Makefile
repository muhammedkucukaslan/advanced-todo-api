swagger: 
	swag fmt
	swag init -g ./infrastructure/fiber/router.go
dev:
	docker-compose up --build 
run: 
	go run ./cmd/main.go

unit-test:
	go test -v ./tests/unit/... 

integration-test:
	go test -v  ./tests/integration/...

httptest:
	go test -v ./tests/httptest/...