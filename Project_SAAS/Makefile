SHELL := /bin/zsh

.PHONY: infra down fmt lint

infra:
	docker compose -f deployments/docker-compose.yml up -d

down:
	docker compose -f deployments/docker-compose.yml down

fmt:
	find . -name '*.go' -print0 | xargs -0 gofmt -w
