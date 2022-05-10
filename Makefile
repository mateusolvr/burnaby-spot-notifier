.PHONY: up
up: ##@development Build and start development environment in background.
	docker-compose up --build -d

.PHONY: shell
shell: ##@development Start a shell session within the container.
	docker-compose run --rm go /bin/sh
	
lint_version ?= v1.40-alpine
.PHONY: lint
lint: ##@development Runs static analysis code.
	docker run --rm \
		-v $(shell pwd):/app \
		-w /app \
		golangci/golangci-lint:$(lint_version) \
		golangci-lint run --timeout 3m

.PHONY: stop
stop: ##@development Stop development environment and remove containers.
	docker-compose down -v --remove-orphans
