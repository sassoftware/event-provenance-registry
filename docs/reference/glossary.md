# Glossary

## Terms

### Event Type

### ADR

An architecture decision record (ADR) is a document that captures an important
architectural decision made along with its context and consequences. Relevant
acronyms:

- AD: architecture decision
- ADL: architecture decision log
- ADR: architecture decision record
- AKM: architecture knowledge management
- ASR: architecturally-significant requirement

### Consumers

Kafka message consumers. Event Provenance Registry (EPR) components that listen
for messages of specific types on the Kafka message bus.

### Event Receiver

The representation of one or more event_types within a event_receiver_group.
Event Receivers are unlocked by events, which track the completion of each
event_receiver's requirements. Event Receivers are reusable lists of
requirements. If and when all the requirements (that is, events) for a
event_receiver all have a "successful" status, the DU can progress through the
event_receiver.

### Event Provenance Registry

Event Provenance Registry (EPR) service that manages and stores events and
tracks event_receivers and event_receiver_groups. Provides an API that let you
create event_receivers and event_receiver_groups. Query the Event Provenance
Registry (EPR) to get identifying information about event_receivers and events.

### graphql

[GraphQL](https://graphql.org/) is a query language for APIs and a runtime for
fulfilling those queries with your existing data. GraphQL provides a complete
and understandable description of the data in your API, gives clients the power
to ask for exactly what they need and nothing more, makes it easier to evolve
APIs over time, and enables powerful developer tools.

### Kafka

An Apache distributed messaging system. Apache Kafka provides fast, highly
scalable and redundant messaging through a publish-subscribe model. Kafka
enables applications to add messages to _topics_. EPR uses a Kafka bus for
messaging. For more information, see the
[Kafka documentation](http://kafka.apache.org/documentation.html).

### Message Bus

Messaging is based on Apache [Kafka](#kafka).

### Success

A Boolean set on a event to determine whether it successful or not.

### NVRPP

Identifier of a DU for purposes of events. Stands for DU **N**ame, **V**ersion,
**R**elease, **P**latform ID, **P**ackage.

For example: `foo-1.0.0-x64-oci-linux-2-docker`

- Name: `foo`
- Version: `1.0.0`
- Release: `20201111.1605122728788`
- Platform ID: `x64-oci-linux-2`
- Package: `docker`

### Partitions

Kafka topic partitions.

### Producers

Kafka Message producers.

### Event

A record of the results of an event_types that occurs in a event_receiver_group.
Takes the form of a CDEvent that contains all the information required to
reproduce the event_type from which it was generated.

Events are created to represent a set of event_types that are performed on a DU.
The events for each event_receiver_group must contain enough information to
recreate the event_receiver_group and to trace the DU version and
[event_receiver](#Event Receiver) that they represent.

### Event Receiver Group

A list of event_receivers. When all event_receivers in a event_receiver_group
have successful events (with "success": true status), the event_receiver_group
publishes a message that specifies an event_type for the event_receiver_group.

### Topics

Category names for Kafka messages.

Examples:

- `epr.dev.events` (Dev)
- `epr.test.events` (Test)
- `epr.prod.events` (Prod)

### Universally Unique Lexicographically Sortable Identifier (ULID)

All IDs in EPR are ULIDs. They are compatible with UUIDs (Universal Unique
Identifier), but they contain a datestamp, are case-insensitive, are
lexicographically sortable, and have fewer characters than UUIDs.

For more information, see [the ULID spec](https://github.com/ulid/spec).
