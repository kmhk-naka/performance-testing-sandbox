#!/bin/bash
export HOST_UID=$(id -u)
export HOST_GID=$(id -g)
docker compose run --rm --service-ports k6-sql "$@"
