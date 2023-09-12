infra-up:
	@docker compose -f .development/docker-compose.yml up -d

infra-down:
	@docker compose -f .development/docker-compose.yml down

test:
	@go test -cover ./...

integration-test:
	@go test -v -count=1 -run ^TestIntegration -tags=integration ./...
