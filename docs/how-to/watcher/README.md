# Watcher

Listed below are instructions on starting a RedPanda instance, setting up
required topics, and starting up a sample watcher to consume records from the
newly created topic.

Watcher code is using the SDK built in to our library.

## Start Redpanda

In a new terminal run the command below. This will start the container up as a
daemon, so be sure to stop the container when you are finished with development.

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

## Create topic

Create topic. Only needed for initial setup

```bash
docker exec -it redpanda-1 \
    rpk topic create epr.dev.events --brokers=localhost:9092
```

## Begin consuming

```bash
cd docs/how-to/watcher
go run main.go
```

You should see a log stating that we have begin consuming records

## Produce message

In a second terminal run the command below

```bash
docker exec -it redpanda-1 \
    rpk topic produce epr.dev.events --brokers=localhost:9092
```

Type a message you would like to send, then press Ctrl+C

```bash
match Ctrl+C
```

## Receive message

You should now see a message like the one below.

```bash
2023/09/01 22:11:19 I received a task with value 'match'
```

**Note**: the matcher being run is looking for kafka messages with the value
`match`. All other messages will be ignored.
