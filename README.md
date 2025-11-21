# Keep Your House Clean API

API em Go para gerenciamento de tarefas domésticas seguindo Clean Architecture e Standard Go Project Layout.

## Stack Tecnológica

- **Linguagem**: Go 1.21
- **Router**: go-chi/chi
- **Database**: PostgreSQL (database/sql nativo) - Suporta Supabase
- **Auth**: golang-jwt/jwt
- **Arquitetura**: Clean Architecture

## Estrutura do Projeto

```
api/
├── cmd/
│   └── api/          # Entry point da aplicação
├── internal/
│   ├── domain/       # Entidades e interfaces de repositório
│   ├── task/         # Service e Handler de tarefas
│   └── platform/     # Infraestrutura (DB, middleware)
migrations/            # Scripts SQL de migração
```

## Pré-requisitos

- Docker e Docker Compose instalados
- Make (opcional, para usar os comandos do Makefile)

## Configuração e Execução

### 1. Inicializar dependências Go

Primeiro, execute o `go mod tidy` usando o Makefile:

```bash
make init-deps
```

Ou manualmente com Docker:

```bash
docker run --rm -v $(pwd):/app -w /app golang:1.21-alpine go mod tidy
```

Ou se preferir executar localmente (requer Go instalado):

```bash
go mod tidy
```

### 2. Executar em modo desenvolvimento

O modo desenvolvimento usa hot-reload com Air:

```bash
make dev-up
```

Ou manualmente:

```bash
docker-compose -f docker-compose.dev.yml up -d
```

### 3. Executar em modo produção

```bash
make build
make up
```

Ou manualmente:

```bash
docker-compose build
docker-compose up -d
```

### 4. Ver logs

```bash
make logs
```

Ou:

```bash
docker-compose logs -f api
```

### 5. Parar os containers

```bash
make down        # Produção
make dev-down    # Desenvolvimento
```

## Variáveis de Ambiente

### Opção 1: DATABASE_URL (Recomendado para Render)

O Render fornece automaticamente `DATABASE_URL` quando você cria um PostgreSQL Database. Use esta variável:

```
DATABASE_URL=postgresql://user:password@host:5432/dbname?sslmode=require
```

### Opção 2: Variáveis Individuais (Para testes locais)

Se preferir usar variáveis individuais:

- `DB_HOST`: Host do banco de dados (padrão: localhost)
- `DB_USER`: Usuário do banco (padrão: postgres)
- `DB_PASSWORD`: Senha do banco (padrão: postgres)
- `DB_NAME`: Nome do banco (padrão: keep_your_house_clean)
- `DB_SSLMODE`: Modo SSL (padrão: disable, use `require` para Render)

### Outras Variáveis

- `JWT_SECRET`: Chave secreta para JWT (obrigatório em produção)
- `PORT`: Porta da API (padrão: 8080)

**Nota**: Se `DATABASE_URL` estiver definida, ela será usada. Caso contrário, usa as variáveis individuais.

## Endpoints da API

### Autenticação

Todas as rotas requerem autenticação via JWT Bearer token no header `Authorization`.

### Tarefas

- `GET /tasks` - Lista todas as tarefas
- `GET /tasks/{id}` - Busca uma tarefa por ID
- `POST /tasks` - Cria uma nova tarefa
- `PUT /tasks/{id}` - Atualiza uma tarefa
- `DELETE /tasks/{id}` - Remove uma tarefa (soft delete)

## Banco de Dados

O PostgreSQL é inicializado automaticamente com Docker Compose. As migrações em `migrations/` são executadas automaticamente. O banco de dados está disponível na porta `5432`.

## Desenvolvimento

Para desenvolvimento com hot-reload, use:

```bash
make dev-up
```

O Air monitora mudanças nos arquivos `.go` e recarrega automaticamente a aplicação.

## Comandos Make Disponíveis

- `make init-deps` - Inicializa dependências Go (go mod tidy)
- `make build` - Builda a imagem Docker
- `make up` - Inicia os containers (produção)
- `make down` - Para os containers (produção)
- `make dev-up` - Inicia os containers (desenvolvimento)
- `make dev-down` - Para os containers (desenvolvimento)
- `make logs` - Mostra os logs da API
- `make clean` - Remove containers e volumes
- `make test` - Executa os testes (requer container rodando com `make dev-up`)
- `make test-coverage` - Executa testes com relatório de cobertura

## Executando Testes

Para executar os testes, primeiro inicie o ambiente de desenvolvimento:

```bash
make dev-up
```

Depois execute os testes:

```bash
make test
```

Ou com cobertura:

```bash
make test-coverage
```

Os testes são executados dentro do container, então não é necessário ter Go instalado localmente.
