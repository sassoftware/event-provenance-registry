# graph

A list of query and mutations are available through the
GraphQL UI. Listed below is an example of running both
a query and a mutation through the UI input fields.

## Query

This query is only returning a subset of the available fields

```graphql
{
  event(id: "1234") {
    name
    version
    description
  }
}
```

## Mutation

Sample Graph QL Queries
Mutation

```graphql
mutation {
  create_event_receiver(
    event_receiver: {
      name: "grant", 
      version: "1.0.0", 
      type: "some-action", 
      description: "a fake event reciever", 
      schema: "{\"name\": \"value\"}"
    }
  )
}
```
