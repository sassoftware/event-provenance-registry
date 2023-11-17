# Single node deployment

Instructions taken from
<https://docs.redpanda.com/22.1/quickstart/quick-start-docker/>

Redpanda can be started up in a single Docker container. The command below will
start it up and continue running in the background. Be sure to stop the
container when you are finished with development.

```bash
docker run -d --pull=always --name=redpanda-1 --rm \
    -p 8081:8081 \
    -p 8082:8082 \
    -p 9092:9092 \
    -p 9644:9644 \
    docker.redpanda.com/redpandadata/redpanda:latest \
    redpanda start \
    --overprovisioned \
    --smp 1  \
    --memory 1G \
    --reserve-memory 0M \
    --node-id 0 \
    --check=false
```

## Create a topic

```bash
 docker exec -it redpanda-1 \
 rpk topic create epr.dev.events --brokers=localhost:9092
```

## Produce a message

```bash
 docker exec -it redpanda-1 \
 rpk topic produce epr.dev.events --brokers=localhost:9092
```

Type text into the topic and press Ctrl + D to separate between messages.

Press Ctrl + C to exit the produce command.

## Consume a message

```bash
 docker exec -it redpanda-1 \
 rpk topic consume epr.dev.events --brokers=localhost:9092
```

Each message is shown with its metadata, like this:

```json
{
  "message": "How do you stream with Redpanda?\n",
  "partition": 0,
  "offset": 1,
  "timestamp": "2021-02-10T15:52:35.251+02:00"
}
```
