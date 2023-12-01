FROM mcr.microsoft.com/devcontainers/go:1-1.21-bookworm

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends postgresql-client
