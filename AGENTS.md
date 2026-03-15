# Repository Guidelines

## Project Structure & Module Organization
- `docker-compose.yml`: shared runtime stack (Go API, MySQL, Prometheus, Grafana, exporters).
- `api-server/`: Go REST API. Main entry is `main.go`; domain split into `handler/`, `repository/`, `model/`, and `seed/`.
- Tool-specific load tests live in `k6/`, `locust/`, `gatling/`, and `artillery/`.
- Observability assets are under `monitoring/` (Prometheus config and Grafana provisioning/dashboards).
- Database bootstrap SQL is in `mysql/init.sql`.

## Build, Test, and Development Commands
- `docker compose up -d`: start shared API/DB/monitoring stack.
- `bash ./test_api.sh`: quick endpoint smoke test (`/health`, CRUD, confirm flow).
- `cd k6 && ./run.sh`: run HTTP load test and generate `k6/results/report.html`.
- `cd k6 && ./run-sql.sh`: run xk6 SQL direct-load test.
- `cd locust && docker compose up -d`: start Locust UI at `http://localhost:8089`.
- `cd gatling && docker compose run --rm gatling`: run Gatling simulation.
- `cd artillery && docker compose run --rm artillery`: run Artillery and generate HTML report.
- `docker compose down` (or `docker compose down -v` for full reset): stop and clean up.

## Coding Style & Naming Conventions
- Go code must be `gofmt`-clean; keep package names lowercase and exported symbols in `CamelCase`.
- Preserve existing JSON/API naming: `snake_case` fields (for example `product_name`, `confirmation_token`).
- Keep endpoint-specific load-test names aligned across tools (`Get Order`, `Create Order`, etc.) for easier comparison.
- Indentation by language: Go tabs (via `gofmt`), YAML 2 spaces, Python 4 spaces.

## Testing Guidelines
- There are currently no committed Go unit tests; new API logic should add `_test.go` files in the same package.
- Before opening a PR, run `bash ./test_api.sh` and at least one load scenario relevant to your change.
- For performance-related changes, include key metrics (throughput, p95 latency, error rate) and reference report paths.

## Commit & Pull Request Guidelines
- Follow the repository’s commit pattern: `<scope>: <short summary>` (example: `k6: xk6-sqlのサンプル追加`).
- Keep commits focused by area (`k6`, `grafana`, `api-server`, etc.).
- PRs should include: purpose, changed directories, commands run, and evidence (report file path or dashboard screenshot) for load/monitoring changes.
