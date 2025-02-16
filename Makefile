run-build:
	docker compose --env-file=docker.env up -d --build

stop:
	docker-compose down

e2e-stand:
	docker compose -f=./docker-compose.e2e.yaml --env-file=docker.e2e.env up -d --build
	#go run ./cmd/app --env-path=local.e2e.env - запускать отдельно (пробовал в разных окружениях, где-то работает полный скрипт, а где-то нет)

test:
	go test -v ./...

test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	rm coverage.out

