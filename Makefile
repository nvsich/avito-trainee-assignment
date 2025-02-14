run:
	docker-compose up -d

run-build:
	docker-compose -d --build

stop:
	docker-compose down

test:
	go test -v ./...