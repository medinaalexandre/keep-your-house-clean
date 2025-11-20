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

### Opção 1: Connection String (Recomendado para Supabase)

- `DATABASE_URL`: Connection string completa do PostgreSQL
  - Exemplo: `postgresql://postgres:[PASSWORD]@db.[PROJECT_REF].supabase.co:5432/postgres`

### Opção 2: Variáveis Individuais do Supabase

- `SUPABASE_DB_HOST`: Host do Supabase (ex: `db.xxxxx.supabase.co`)
- `SUPABASE_DB_PORT`: Porta (padrão: `5432`)
- `SUPABASE_DB_USER`: Usuário (padrão: `postgres`)
- `SUPABASE_DB_PASSWORD`: Senha do banco
- `SUPABASE_DB_NAME`: Nome do banco (padrão: `postgres`)
- `SUPABASE_DB_SSLMODE`: Modo SSL (padrão: `require`)

### Opção 3: Variáveis Tradicionais (Compatibilidade)

- `DB_HOST`: Host do banco de dados (padrão: localhost)
- `DB_USER`: Usuário do banco (padrão: postgres)
- `DB_PASSWORD`: Senha do banco (padrão: postgres)
- `DB_NAME`: Nome do banco (padrão: keep_your_house_clean)
- `DB_SSLMODE`: Modo SSL (padrão: disable)

### Outras Variáveis

- `JWT_SECRET`: Chave secreta para JWT (obrigatório em produção)
- `PORT`: Porta da API (padrão: 8080)

**Nota**: O código detecta automaticamente qual método usar. Prioridade: `DATABASE_URL` > `SUPABASE_DB_*` > `DB_*`

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

### Desenvolvimento Local

O PostgreSQL é inicializado automaticamente com Docker Compose. As migrações em `migrations/` são executadas automaticamente. O banco de dados está disponível na porta `5432`.

### Produção com Supabase

O projeto suporta Supabase como banco de dados PostgreSQL. Veja o guia completo em [SUPABASE_SETUP.md](./SUPABASE_SETUP.md).

**Configuração rápida:**
1. Crie um projeto no Supabase (https://supabase.com)
2. Configure a variável `DATABASE_URL` com a connection string
3. Execute as migrações via SQL Editor do Supabase
4. Pronto! A API se conecta automaticamente ao Supabase

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
