.PHONY: build migrate run run-bg stop run-dev

build:
	go build -o bin/app main.go
	go build -o bin/migrate scripts/migration.go

run: build
	./bin/app

run-bg: build
	@mkdir -p logs
	@nohup ./bin/app > logs/app.log 2>&1 & echo $$! > bin/app.pid
	@echo "App is running in background. PID: `cat bin/app.pid`"

stop:
	@echo "Stopping app..."
	@if [ -f bin/app.pid ]; then \
		kill `cat bin/app.pid` && rm bin/app.pid && echo "Stopped."; \
	else \
		echo "No app.pid file found. Is the app running?"; \
	fi

migrate: build
	./bin/migrate

run-dev:
	go run main.go


docker-migrate:
	./migrate


