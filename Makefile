deps:
	@docker-compose -f dev/docker-compose.yaml up -d

deps-down:
	@docker-compose -f dev/docker-compose.yaml down

run-api:
	go run . api
