# Deploying Redpanda

## Overview

An event streaming platform conforming to the Apache Kafka API is required to
run the `async-event-infrastructure` server. This how-to is intended to guide a
user through deploying one such platform, Redpanda, on their local machine.

## Prerequisites

* Docker

## Installation

First copy the file [docker-compose.yaml](./docker-compose.yaml) to your
machine.

Navigate to the directory where you copied the file in a command prompt. Then
use Docker Compose to deploy Redpanda.

```bash
docker compose up -d
```

A single Redpanda broker is run, configured without authentication and
accessible from `localhost:19092`. A web console is also hosted and can be
accessed in your browser at http://localhost:8080/overview.

To later tear down the Redpanda deployment, run the following command from the
same directory as before with access to your copy
of [docker-compose.yaml](./docker-compose.yaml).

```bash
docker compose down -v
```

## Documentation

- [What is Redpanda](https://redpanda.com/what-is-redpanda)
- [Redpanda documentation](https://docs.redpanda.com/docs/home/)
- [Install the Redpanda CLI](https://docs.redpanda.com/docs/get-started/rpk-install/)
