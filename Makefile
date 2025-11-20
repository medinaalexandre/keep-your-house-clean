.PHONY: help build up down dev-up dev-down logs clean test test-coverage init-deps seed

help:
	@echo "Comandos disponíveis:"
	@echo "  make init-deps  - Inicializa dependências Go (go mod tidy)"
	@echo "  make build      - Builda a imagem Docker"
	@echo "  make up         - Inicia os containers (produção)"
	@echo "  make down       - Para os containers (produção)"
	@echo "  make dev-up     - Inicia os containers (desenvolvimento)"
	@echo "  make dev-down   - Para os containers (desenvolvimento)"
	@echo "  make logs       - Mostra os logs da API"
	@echo "  make clean      - Remove containers e volumes"
	@echo "  make test       - Executa os testes (requer container rodando)"
	@echo "  make test-coverage - Executa testes com cobertura (requer container rodando)"
	@echo "  make seed       - Popula o banco de dados com dados de exemplo (requer container rodando)"

init-deps:
	@./scripts/init-deps.sh

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

dev-up:
	docker-compose -f docker-compose.dev.yml up -d

dev-down:
	docker-compose -f docker-compose.dev.yml down

logs:
	docker-compose logs -f api

clean:
	docker-compose down -v
	docker-compose -f docker-compose.dev.yml down -v

test:
	@if docker ps --format '{{.Names}}' | grep -q "^keep-your-house-clean-dev$$"; then \
		docker-compose -f docker-compose.dev.yml exec api go test -v ./...; \
	elif docker ps --format '{{.Names}}' | grep -q "^keep-your-house-clean$$"; then \
		docker-compose exec api go test -v ./...; \
	else \
		echo "Nenhum container rodando. Execute 'make dev-up' primeiro."; \
		exit 1; \
	fi

test-coverage:
	@if docker ps --format '{{.Names}}' | grep -q "^keep-your-house-clean-dev$$"; then \
		docker-compose -f docker-compose.dev.yml exec api sh -c "go test -v -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html"; \
	elif docker ps --format '{{.Names}}' | grep -q "^keep-your-house-clean$$"; then \
		docker-compose exec api sh -c "go test -v -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html"; \
	else \
		echo "Nenhum container rodando. Execute 'make dev-up' primeiro."; \
		exit 1; \
	fi

seed:
	@if docker ps --format '{{.Names}}' | grep -q "^keep-your-house-clean-dev$$"; then \
		docker-compose -f docker-compose.dev.yml exec api go run cmd/seeder/main.go; \
	elif docker ps --format '{{.Names}}' | grep -q "^keep-your-house-clean$$"; then \
		docker-compose exec api go run cmd/seeder/main.go; \
	else \
		echo "Nenhum container rodando. Execute 'make dev-up' primeiro."; \
		exit 1; \
	fi
