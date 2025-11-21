.PHONY: help build build-prod up down dev-up dev-down logs clean test test-coverage init-deps seed run-prod

help:
	@echo "Comandos disponíveis:"
	@echo "  make init-deps     - Inicializa dependências Go (go mod tidy)"
	@echo "  make build         - Builda a imagem Docker (desenvolvimento)"
	@echo "  make build-prod    - Builda a imagem Docker de produção (Dockerfile.prod)"
	@echo "  make run-prod      - Executa a imagem de produção localmente"
	@echo "  make up            - Inicia os containers (produção)"
	@echo "  make down          - Para os containers (produção)"
	@echo "  make dev-up        - Inicia os containers (desenvolvimento)"
	@echo "  make dev-down      - Para os containers (desenvolvimento)"
	@echo "  make logs          - Mostra os logs da API"
	@echo "  make clean         - Remove containers e volumes"
	@echo "  make test          - Executa os testes (requer container rodando)"
	@echo "  make test-coverage - Executa testes com cobertura (requer container rodando)"
	@echo "  make seed          - Popula o banco de dados com dados de exemplo (requer container rodando)"

init-deps:
	@./scripts/init-deps.sh

build:
	docker-compose build

build-prod:
	@echo "Building production image..."
	docker build -f Dockerfile.prod -t keep-your-house-clean:prod .
	@echo "Build concluído! Imagem: keep-your-house-clean:prod"

run-prod:
	@echo "Executando imagem de produção..."
	@echo "Certifique-se de configurar as variáveis de ambiente antes!"
	@echo "Exemplo:"
	@echo "  export DATABASE_URL='postgresql://...'"
	@echo "  export JWT_SECRET='sua-chave'"
	@echo "  make run-prod"
	@echo ""
	docker run -p 8080:8080 \
		-e DATABASE_URL="$${DATABASE_URL}" \
		-e DB_HOST="$${DB_HOST}" \
		-e DB_PORT="$${DB_PORT}" \
		-e DB_USER="$${DB_USER}" \
		-e DB_PASSWORD="$${DB_PASSWORD}" \
		-e DB_NAME="$${DB_NAME}" \
		-e DB_SSLMODE="$${DB_SSLMODE}" \
		-e JWT_SECRET="$${JWT_SECRET}" \
		keep-your-house-clean:prod

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
