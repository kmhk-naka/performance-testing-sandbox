#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPT_DIR}"

MYSQL_HOST="${MYSQL_HOST:-mysql}"
MYSQL_PORT="${MYSQL_PORT:-3306}"
MYSQL_USER="${MYSQL_USER:-app}"
MYSQL_PASSWORD="${MYSQL_PASSWORD:-password}"
MYSQL_DB="${MYSQL_DB:-orders_db}"

TEST_NAME="${TEST_NAME:-oltp_read_write}"
TABLES="${TABLES:-4}"
TABLE_SIZE="${TABLE_SIZE:-10000}"
THREADS="${THREADS:-16}"
TIME_SECONDS="${TIME_SECONDS:-60}"
REPORT_INTERVAL="${REPORT_INTERVAL:-2}"

phase="${1:-all}"
if [[ $# -gt 0 ]]; then
  shift
fi
EXTRA_ARGS=("$@")

common_args=(
  "--db-driver=mysql"
  "--mysql-host=${MYSQL_HOST}"
  "--mysql-port=${MYSQL_PORT}"
  "--mysql-user=${MYSQL_USER}"
  "--mysql-password=${MYSQL_PASSWORD}"
  "--mysql-db=${MYSQL_DB}"
)

prepare_args=(
  "${TEST_NAME}"
  "${common_args[@]}"
  "--tables=${TABLES}"
  "--table-size=${TABLE_SIZE}"
  "${EXTRA_ARGS[@]}"
  prepare
)

run_args=(
  "${TEST_NAME}"
  "${common_args[@]}"
  "--tables=${TABLES}"
  "--table-size=${TABLE_SIZE}"
  "--threads=${THREADS}"
  "--time=${TIME_SECONDS}"
  "--report-interval=${REPORT_INTERVAL}"
  "${EXTRA_ARGS[@]}"
  run
)

cleanup_args=(
  "${TEST_NAME}"
  "${common_args[@]}"
  "--tables=${TABLES}"
  "${EXTRA_ARGS[@]}"
  cleanup
)

run_sysbench() {
  docker compose run --rm sysbench "$@"
}

run_benchmark() {
  mkdir -p results
  local timestamp
  timestamp="$(date +%Y%m%d-%H%M%S)"
  local log_file
  log_file="results/${TEST_NAME}-${timestamp}.log"

  run_sysbench "${run_args[@]}" | tee "${log_file}"
  echo "Result log: ${log_file}"
}

case "${phase}" in
  prepare)
    run_sysbench "${prepare_args[@]}"
    ;;
  run)
    run_benchmark
    ;;
  cleanup)
    run_sysbench "${cleanup_args[@]}"
    ;;
  all)
    run_sysbench "${prepare_args[@]}"
    run_benchmark
    ;;
  *)
    cat <<'USAGE'
Usage:
  ./run.sh [prepare|run|cleanup|all] [extra sysbench options...]

Examples:
  ./run.sh all
  ./run.sh prepare
  ./run.sh run --rand-type=uniform
  TEST_NAME=oltp_read_only THREADS=32 TIME_SECONDS=120 ./run.sh run
USAGE
    exit 1
    ;;
esac
