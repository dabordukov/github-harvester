# GitHub Harvester

A microservice project for retrieving GitHub repository information and managing subscriptions to them.

## Services

- `api` — HTTP gateway. Provides endpoints for subscription management and repository information retrieval.
- `subscriber` — stores user subscriptions in its own PostgreSQL database.
- `processor` — the central access point for repository data. Acts as a cache and background update orchestrator. Stores collected information in its own PostgreSQL database. On incoming requests, it first checks the database (cache); if data is missing, it publishes a collection task to Kafka.
- `collector` — asynchronous worker. Consumes data collection tasks from Kafka, calls the GitHub API, and publishes results back to Kafka. Has no direct synchronous interaction with `processor`.
- `kafka` — message broker providing asynchronous communication between `processor` and `collector`.
- `postgres` — `subscriber` and `processor` each have their own isolated databases/schemas.

### Single Request Processing (Cache-Aside):
1. Client requests repository information via `api` from `processor`.
2. `Processor` looks up data in its PostgreSQL database (reads from cache).
3. If the repository is not found, `processor` publishes a collection task message to a Kafka topic.
4. `Collector` consumes the message from Kafka, makes a request to the GitHub API, and publishes the collected data back to Kafka.
5. `Processor` reads the result from Kafka, persists it to its database, and returns the response to the client.

### Background Subscription Sync:
- Every 15 seconds, `collector` independently initiates the cache update process.
- It requests the current list of subscriptions from `subscriber`.
- Based on the received list, `collector` generates update tasks and sends them to Kafka, ensuring regular data updates for subscriptions in the `processor` database without direct user involvement.

## API Endpoints

- `GET /api/ping`
- `GET /api/repositories/info?url=https://github.com/{owner}/{repo}`
- `POST /api/subscriptions`
- `DELETE /api/subscriptions/{owner}/{repo}`
- `GET /api/subscriptions`
- `GET /api/subscriptions/info`

## Running Locally with Docker Compose

From the repository root:

```bash
docker compose up --build
```

After startup, the services are available at:

- Swagger UI: `http://localhost:28080/swagger/`
- API: `http://localhost:28080`

## Request Examples

Get information about a specific repository:

```bash
curl -X GET "http://localhost:28080/api/repositories/info?url=https://github.com/golang/go"
```

Subscribe to a repository:

```bash
curl -X POST http://localhost:28080/api/subscriptions \
  -H "Content-Type: application/json" \
  -d '{"owner":"golang","repo_name":"go"}'
```

Get the list of subscriptions:

```bash
curl -X GET http://localhost:28080/api/subscriptions
```

Get aggregated information on subscriptions:

```bash
curl -X GET http://localhost:28080/api/subscriptions/info
```

Unsubscribe:

```bash
curl -X DELETE http://localhost:28080/api/subscriptions/golang/go
```
