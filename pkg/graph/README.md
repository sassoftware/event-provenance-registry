# graph

This readme is  for instructions on working with gqlgen

To update generated code from schema updates, update the
graphql schema at [pkg/graph/schema.graphqls](./schema.graphqls),
then run the command below.

```bash
go run github.com/99designs/gqlgen generate
```

Listed below is an example graphql query that can be used to create
a new event reciever

Mutation

```graphql
mutation CreateEventReceiver($input: EventReceiverInput!) {
  create_event_receiver(input: $input) {
    name
    type
    version
    description
    enabled
  }
}
```

Variables

```json
{
  "input": {
    "name": "name",
    "type": "type",
    "version": "version",
    "description": "description",
    "enabled": true
  }
}
```
