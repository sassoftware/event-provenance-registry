FROM mcr.microsoft.com/devcontainers/go:1-1.21-bookworm

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends postgresql-client

RUN curl -LO https://github.com/redpanda-data/redpanda/releases/download/v23.2.7/rpk-linux-amd64.zip \
    && unzip rpk-linux-amd64.zip \
    && install rpk -t /usr/local/bin \
    && rm rpk-linux-amd64.zip
