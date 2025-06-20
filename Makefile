# Config
PROJECT_NAME = fiber-mvc-api
DOCKER_DEV_COMPOSE = docker/docker-compose.dev.yml
DOCKER_PROD_COMPOSE = docker/docker-compose.yml

# Dev containers
dev-up:
	docker compose -f $(DOCKER_DEV_COMPOSE) up --build

dev-up-d:
	docker compose -f $(DOCKER_DEV_COMPOSE) up --build -d

dev-down:
	docker compose -f $(DOCKER_DEV_COMPOSE) down

dev-down-v:
	docker compose -f $(DOCKER_DEV_COMPOSE) down -v --remove-orphans

dev-logs:
	docker compose -f $(DOCKER_DEV_COMPOSE) logs -f

dev-sh:
	docker exec -it $(PROJECT_NAME)-dev sh

# Prod containers
prod-up:
	docker compose -f $(DOCKER_PROD_COMPOSE) up --build -d

prod-down:
	docker compose -f $(DOCKER_PROD_COMPOSE) down

prod-down-v:
	docker compose -f $(DOCKER_PROD_COMPOSE) down -v --remove-orphans

prod-logs:
	docker compose -f $(DOCKER_PROD_COMPOSE) logs -f

prod-sh:
	docker exec -it $(PROJECT_NAME) sh

# Clean dangling stuff
prune:
	docker system prune -f --volumes

# Go commands
run-main:
	go run main.go

run-build:
	go run build.go

build:
	go build -o bin/app main.go

clean:
	rm -rf bin/

# Show container status
ps:
	docker ps --filter name=$(PROJECT_NAME)