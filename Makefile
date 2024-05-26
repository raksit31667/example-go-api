.PHONY: setup-pre-commit
setup-pre-commit:
	@echo "Setting up pre-commit..."
	@./scripts/setup-pre-commit.sh

.PHONY: setup-local-db
setup-local-db:
	@echo "Starting local database..."
	docker-compose -f docker-compose.yml down && \
	docker-compose -f docker-compose.yml up --build --force-recreate -d

.PHONY: run-local
run-local:
	@echo "Running the application locally..."
	go run main.go

.PHONY: test
test:
	@echo "Running unit tests..."
	go test -v ./...

.PHONY: test-coverage
test-coverage:
	@echo "Running unit tests and generating test coverage report..."
	go test -v ./... -coverprofile=coverage.out

.PHONY: integration-test
integration-test:
	@echo "Running integration tests..."
	docker-compose -f docker-compose.it.test.yml down && \
	docker-compose -f docker-compose.it.test.yml up --build --force-recreate --abort-on-container-exit --exit-code-from integration-test
