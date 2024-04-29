# Glossary

## Terms

### Event Provenance Registry

Event Provenance Registry (EPR) is a service designed to efficiently manage and
store events, while also tracking event_receivers and event_receiver_groups.
Offering a robust API, EPR empowers users to seamlessly create event_receivers
and event_receiver_groups. Furthermore, it facilitates querying the Event
Provenance Registry, enabling users to obtain crucial identifying information
about event_receivers and events. Through these capabilities, EPR plays a
central role in enhancing the visibility, traceability, and management of events
within a system.

### Events

Events represent the core data entities within EPR. An event is simply a record
of something changing state. These entities encapsulate information about tasks
executed in the pipeline. Events are created to represent a set of event_types
that are performed on a unit. Events act as the foundational building blocks of
the provenance tracking system.

### Event Receivers

Event receivers are structures within EPR designed to capture and catagorize
events. They act as recipients of specific events, defining the scope and nature
of events they are interested in. These receivers streamline the categorization
and handling of events based on predetermined criteria. Event Receivers are
reusable lists of requirements. When an event receiver receives an event, it
publishes a message that specifies the type of event.

### Event Receiver Groups

Event receiver groups are collections of event receivers. Event receiver groups
provide a higher-level organizational structure within EPR. By grouping multiple
event receivers, users can efficiently manage and coordinate the processing of
events across various categories or functional units. This grouping mechanism
enhances the scalability and modularity of event handling.

When all event_receivers in an event_receiver_group have successful events (with
"success": true status), the event_receiver_group publishes a message that
specifies an event_type for the event_receiver_group.

An event_receiver_group can be enabled or disabled (with "enabled": true|false).
When enabled, the event_receiver_group acts as previously described. When
disabled, messages for the event_receiver_group will not be published. Messages
will still be published for the underlying event_receivers should they receive
successful events.

### NVRPP

A crucial aspect of event identification in EPR is the NVRPP identifier,
standing for Name, Version, Release, Platform ID, and Package. These five keys
collectively form a unique identifier for an event within the pipeline. The
NVRPP serves as a tracking mechanism, allowing for the association of events
with specific units in the pipeline. This identifier is particularly significant
for receipts, enabling a comprehensive understanding of events within a specific
context.

For example: `foo-1.0.0-x64-oci-linux-2-docker`

- Name: `foo`
- Version: `1.0.0`
- Release: `20201111.1605122728788`
- Platform ID: `x64-oci-linux-2`
- Package: `docker`

### Success

A Boolean set on a event to determine whether it successful or not.

### graphql

[GraphQL](https://graphql.org/) is a query language for APIs and a runtime for
fulfilling those queries with your existing data. GraphQL provides a complete
and understandable description of the data in your API, gives clients the power
to ask for exactly what they need and nothing more, makes it easier to evolve
APIs over time, and enables powerful developer tools.

### Redpanda

Redpanda is a simple, powerful, and cost-efficient streaming data platform that
is compatible with Kafka® APIs while eliminating Kafka complexity. Redpanda
offers a complete streaming data platform in a single binary, including brokers,
HTTP proxy, and schema registry services. For more information, see the
[Redpanda documentation](https://docs.redpanda.com/current/home).

### Kafka

An Apache distributed messaging system. Apache Kafka provides fast, highly
scalable and redundant messaging through a publish-subscribe model. Kafka
enables applications to add messages to _topics_. EPR uses a Kafka bus for
messaging. For more information, see the
[Kafka documentation](http://kafka.apache.org/documentation.html).

### Message Bus

Messaging is based on Apache [Kafka](#kafka).

### Consumers

Kafka message consumers. Event Provenance Registry (EPR) components that listen
for messages of specific types on the Kafka message bus.

### Producers

Kafka Message producers.

### Topics

Category names for Kafka messages.

Examples:

- `epr.dev.events` (Dev)
- `epr.test.events` (Test)
- `epr.prod.events` (Prod)

### Partitions

Kafka topic partitions.

### Universally Unique Lexicographically Sortable Identifier (ULID)

All IDs in EPR are ULIDs. They are compatible with UUIDs (Universal Unique
Identifier), but they contain a datestamp, are case-insensitive, are
lexicographically sortable, and have fewer characters than UUIDs.

For more information, see [the ULID spec](https://github.com/ulid/spec).

### ADR

An architecture decision record (ADR) is a document that captures an important
architectural decision made along with its context and consequences. Relevant
acronyms:

- AD: architecture decision
- ADL: architecture decision log
- ADR: architecture decision record
- AKM: architecture knowledge management
- ASR: architecturally-significant requirement

# Request Enhancement Proposals (REP)

A Request Enhancement Proposal (REP) is a way to propose, communicate about, and
coordinate new initiatives for the project.
