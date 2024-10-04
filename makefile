.PHONY: default run build test docs clean
# Variables
APP_NAME=FORNECEDORES

# Tasks
default: run-with-docs

run:
	@go run main.go
run-with-docs:
	@swag init --parseDependency --parseInternal
	@go run main.go
build:
	@go build -o $(APP_NAME) main.go
test:
	@go test ./ ...
docs:
	@swag init
clean:
	@rm -f $(APP_NAME)
	@rm -rf ./docs