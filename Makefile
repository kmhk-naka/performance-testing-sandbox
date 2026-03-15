SHELL := /bin/bash

.DEFAULT_GOAL := help

MODE ?= all
ARGS ?=

.PHONY: help \
	api-tools-up api-tools-down db-tools-up db-tools-down \
	k6-up k6-down k6-sql-up k6-sql-down \
	locust-up locust-down \
	gatling-up gatling-down artillery-up artillery-down \
	sysbench-up sysbench-down \
	down down-v ps logs health-api test-api \
	k6 k6-sql gatling artillery sysbench \
	sysbench-prepare sysbench-run sysbench-cleanup

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Scenario Stack:"
	@echo "  api-tools-up      Start API tools stack (api + mysql + monitoring)"
	@echo "  api-tools-down    Stop API tools stack"
	@echo "  db-tools-up       Start DB direct tools stack (mysql + monitoring)"
	@echo "  db-tools-down     Stop DB direct tools stack"
	@echo "  down              Stop current root stack"
	@echo "  down-v            Stop root stack and remove volumes"
	@echo "  ps                Show running services"
	@echo "  logs              Follow logs from all running services"
	@echo ""
	@echo "Tool Up/Down:"
	@echo "  k6-up / k6-down"
	@echo "  k6-sql-up / k6-sql-down"
	@echo "  locust-up / locust-down"
	@echo "  gatling-up / gatling-down"
	@echo "  artillery-up / artillery-down"
	@echo "  sysbench-up / sysbench-down"
	@echo ""
	@echo "Checks:"
	@echo "  health-api        Check API health endpoint"
	@echo "  test-api          Run API smoke test script"
	@echo ""
	@echo "Load Tests (auto-start required stack):"
	@echo "  k6                Run k6 HTTP load test"
	@echo "  k6-sql            Run k6 SQL direct load test"
	@echo "  gatling           Run Gatling simulation"
	@echo "  artillery         Run Artillery load test"
	@echo "  sysbench          Run sysbench (MODE/ARGS configurable)"
	@echo "  locust-up         Start Locust UI with required stack"
	@echo "  locust-down       Stop Locust UI and stack"
	@echo "  sysbench-prepare  Prepare sysbench data"
	@echo "  sysbench-run      Run sysbench benchmark only"
	@echo "  sysbench-cleanup  Cleanup sysbench data"
	@echo ""
	@echo "Examples:"
	@echo "  make k6-up && make k6 && make k6-down"
	@echo "  make locust-up  # open http://localhost:8089"
	@echo "  make locust-down"
	@echo "  make sysbench-up && make sysbench && make sysbench-down"
	@echo "  make k6-sql"
	@echo "  make sysbench-run ARGS='--rand-type=uniform'"

api-tools-up:
	docker compose --profile api --profile monitoring up -d

api-tools-down:
	docker compose down

db-tools-up:
	docker compose --profile db --profile monitoring up -d

db-tools-down:
	docker compose down

down:
	$(MAKE) api-tools-down

down-v:
	docker compose down -v

ps:
	docker compose ps

logs:
	docker compose logs -f

health-api:
	curl -fsS http://localhost:8080/health

test-api:
	bash ./test_api.sh

k6-up: api-tools-up

k6-down: api-tools-down

k6: k6-up
	cd k6 && ./run.sh

k6-sql-up: db-tools-up

k6-sql-down: db-tools-down

k6-sql: k6-sql-up
	cd k6 && ./run-sql.sh

locust-up: api-tools-up
	cd locust && docker compose up -d

locust-down:
	-cd locust && docker compose down
	docker compose down

gatling-up: api-tools-up

gatling-down: api-tools-down

gatling: gatling-up
	cd gatling && docker compose run --rm gatling

artillery-up: api-tools-up

artillery-down: api-tools-down

artillery: artillery-up
	cd artillery && docker compose run --rm artillery

sysbench-up: db-tools-up

sysbench-down: db-tools-down

sysbench: sysbench-up
	cd sysbench && ./run.sh $(MODE) $(ARGS)

sysbench-prepare:
	$(MAKE) sysbench MODE=prepare ARGS="$(ARGS)"

sysbench-run:
	$(MAKE) sysbench MODE=run ARGS="$(ARGS)"

sysbench-cleanup:
	$(MAKE) sysbench MODE=cleanup ARGS="$(ARGS)"
