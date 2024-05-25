.PHONY: test
test:
	@echo "Running unit tests..."
	go test -v ./...

.PHONY: integration-test
integration-test:
	@echo "Running integration tests..."
	docker-compose -f docker-compose.it.test.yml down && \
	docker-compose -f docker-compose.it.test.yml up --build --force-recreate --abort-on-container-exit --exit-code-from integration-test
