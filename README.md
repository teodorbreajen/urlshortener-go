# URL Shortener

API REST para acortar URLs construida en Go con arquitectura hexagonal.

## Características

- Crear URLs cortas
- Redirección rápida
- Estadísticas de uso
- Arquitectura hexagonal (Ports & Adapters)

## Estructura

`
cmd/server/           # Punto de entrada
internal/
  adapter/           # Adaptadores HTTP
  domain/            # Modelo de dominio
  port/              # Interfaces
`

## Uso

\\\ash
go mod download
go run cmd/server/main.go
\\\

## API Endpoints

- POST /urls - Crear URL corta
- GET /{short} - Redireccionar a URL original
- GET /urls/{short}/stats - Ver estadísticas

## Tecnologías

- Go 1.21+
- Gin/Gorilla Mux
- SQLite/PostgreSQL
- Arquitectura Hexagonal