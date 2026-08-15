# GitHub Harvester

Микросервисный проект для получения информации о репозиториях GitHub и управления подписками на них.

## Сервисы

- `api` — HTTP gateway. Отдаёт старый endpoint по URL репозитория и новые endpoints для подписок.
- `subscriber` — хранит подписки в PostgreSQL, проверяет существование репозитория через GitHub REST API.
- `processor` — orchestration слой между `api` и `collector`.
- `collector` — получает данные о репозиториях из GitHub и умеет собирать информацию по всем подпискам.
- `postgres` — база данных для `subscriber`.

## Эндпоинты API

- `GET /api/ping`
- `GET /api/repositories/info?url=https://github.com/{owner}/{repo}`
- `POST /api/subscriptions`
- `DELETE /api/subscriptions/{owner}/{repo}`
- `GET /api/subscriptions`
- `GET /api/subscriptions/info`

## Запуск локально через Docker Compose

Из корня репозитория:

```bash
docker compose up --build
```

После запуска сервисы доступны на:

- API: `http://localhost:28080`
- Subscriber gRPC: `localhost:28081`
- Processor gRPC: `localhost:28082`
- Collector gRPC: `localhost:28083`

## Примеры запросов

Получить информацию о конкретном репозитории:

```bash
curl -X GET "http://localhost:28080/api/repositories/info?url=https://github.com/golang/go"
```

Подписаться на репозиторий:

```bash
curl -X POST http://localhost:28080/api/subscriptions \
  -H "Content-Type: application/json" \
  -d '{"owner":"golang","repo_name":"go"}'
```

Получить список подписок:

```bash
curl -X GET http://localhost:28080/api/subscriptions
```

Получить агрегированную информацию по подпискам:

```bash
curl -X GET http://localhost:28080/api/subscriptions/info
```

Отписаться:

```bash
curl -X DELETE http://localhost:28080/api/subscriptions/golang/go
```

## Локальная разработка без Docker

В каталоге `repo-stat`:

```bash
make protobuf
make sqlc
go test ./...
```

Для `subscriber` нужно поднять PostgreSQL и передать `DATABASE_DSN`.
