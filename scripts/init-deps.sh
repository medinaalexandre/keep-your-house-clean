#!/bin/sh

echo "Inicializando dependências Go..."

docker run --rm \
  -v $(pwd):/app \
  -w /app \
  golang:1.21-alpine \
  sh -c "go mod tidy && go mod download"

echo "Dependências inicializadas com sucesso!"
