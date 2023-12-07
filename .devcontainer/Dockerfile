FROM mcr.microsoft.com/devcontainers/go:1-1.21-bookworm

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends postgresql-client

ARG REDPANDA_VERSION=23.2.17

RUN mkdir -p /tmp/redpanda \
    && cd /tmp/redpanda \
    && curl -LO "https://github.com/redpanda-data/redpanda/releases/download/v${REDPANDA_VERSION}/rpk-linux-amd64.zip" \
    && curl -LO "https://github.com/redpanda-data/redpanda/releases/download/v${REDPANDA_VERSION}/rpk_${REDPANDA_VERSION}_checksums.txt" \
    && sha256sum -c rpk_${REDPANDA_VERSION}_checksums.txt --ignore-missing --strict \
    && unzip rpk-linux-amd64.zip \
    && install rpk -t /usr/local/bin \
    && rm -r /tmp/redpanda
