# GitHub Harvester

Распределённая система для анализа репозиториев GitHub. Состоит из двух микросервисов:

- **Gateway** — HTTP API с gRPC клиентом
- **Collector** — gRPC сервис для сбора данных о GitHub репозиториях

## Архитектура

```
┌─────────────┐      gRPC       ┌─────────────┐
│   Gateway   │ ◄──────────────►│  Collector  │
│  (HTTP API) │                 │ (gRPC svc)  │
└─────────────┘                 └─────────────┘
                                        │
                                        ▼
                                 ┌─────────────┐
                                 │  GitHub API │
                                 └─────────────┘
```

## Структура проекта

- `cmd/collector/` — gRPC микросервис для сбора данных
- `cmd/gateway/` — HTTP API Gateway
- `internal/` — общие пакеты
- `demo/` — демо-версия (cli)

## Запуск

### Переменные окружения

- `COLLECTOR_PORT` — порт gRPC сервера (по умолчанию: 8888)
- `HTTP_PORT` — порт HTTP сервера (по умолчанию: 8080)

### Запуск через Docker Compose

```bash
docker-compose up -d
```

### Ручной запуск

1. Запустить Collector:
```bash
export COLLECTOR_PORT=8888
go run cmd/collector/main.go
```

2. Запустить Gateway:
```bash
export COLLECTOR_ADDR=localhost:8888
export HTTP_PORT=8080
go run cmd/gateway/main.go
```

## API

### Получить информацию о репозитории

```
GET /repo/{owner}/{repo}
```

**Пример:**
```bash
curl -H "Authorization: github_token" http://localhost:8080/repo/golang/go
```

**Ответ:**
```json
{
  "name": "go",
  "owner": "",
  "description": "The Go programming language",
  "forks": 18860,
  "stars": 133046,
  "created_at": "2014-08-19T04:33:40Z",
  "commits_count": 65642
}
```

### Swagger документация

```
GET /swagger/
```

## Demo версия

CLI инструмент:

```bash
export GITHUB_TOKEN=token
go run demo/main.go
```
