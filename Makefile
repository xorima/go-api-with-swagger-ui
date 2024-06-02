.PHONY: test build docs run grafana-run

test: docs
	go test ./...

build: docs test
	go build

docs:
	swag init

run: docs test
	podman-compose up --build -d

restart:
	podman-compose restart
