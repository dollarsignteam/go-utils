infra-up:
	@docker compose -f .development/docker-compose.yml up -d

infra-down:
	@docker compose -f .development/docker-compose.yml down

test:
	@go test -cover ./...

integration-test:
	@go test -v -run ^TestIntegration -tags=integration ./...