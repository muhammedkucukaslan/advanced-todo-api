swagger: 
	swag fmt
	swag init -g ./infrastructure/fiber/router.go
dev:
	docker-compose up --build 
run: 
	go run ./cmd/main.go

test: unit-test integration-test httptest e2e-test

unit-test:
	go test -v ./tests/unit/... 

integration-test:
	go test -v  ./tests/integration/...

httptest:
	go test -v ./tests/httptest/...

e2e-test:
	go test -v ./tests/e2e/...

compose:
	docker-compose up --build