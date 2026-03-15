FROM debian:bookworm-slim

RUN apt-get update \
  && apt-get install -y --no-install-recommends sysbench default-mysql-client ca-certificates \
  && rm -rf /var/lib/apt/lists/*

ENTRYPOINT ["sysbench"]
