.PHONY: up
up: ##@development Build and start development environment in background.
	docker-compose up --build -d

VERSION = 1.1.0
tag: 
	docker tag web-scraper-go-go:latest mateusolvr/personal:web-scraper-burnaby-v$(VERSION)

push:
	docker push mateusolvr/personal:web-scraper-burnaby-v$(VERSION)

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


# clean: ##@dev Remove folder vendor, public and coverage.
# 	rm -rf vendor public coverage
# lint: clean ##@check Run lint on docker.
# 	DOCKER_BUILDKIT=1 \
# 	docker build --progress=plain \
# 		--target=lint \
# 		--file=./Dockerfile .

.PHONY: stop
stop: ##@development Stop development environment and remove containers.
	docker-compose down -v --remove-orphans
