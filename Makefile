
build:
	go build -o build/app main.go

api-docs:
	swag init --output docs

migrate:
	./migrate