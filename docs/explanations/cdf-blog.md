# EPR

Event Provenance Registry is culmination of several years of SAS's effort to
convert from large ship events to CI/CD. We built the first version internally
as a way to facilitate CI/CD in a complex, aging build system. The result
enables SAS
to build, package, scan, promote, and ship thousands of artifacts daily.

## Origin

Around mid 2019, management gave the directive that R&D needed to move faster to
stay competitive. SAS's software release process to that point consisted of two
or three large ship events per year. Updates were slow and painful. To reverse
that trend, we needed a way to shorten the development cycle and deliver
artifacts more quickly.
The challenge was twofold. First, create a system that could allow disparate
pieces of our pipeline to communicate and chain together. Second, help R&D shift
gears from the old software development model to CI/CD.

In a happy, imaginary world somewhere, we could have chained together Github
actions or equivalent into a working pipeline. No need to write EPR at all.
Reality was not so kind. Our source code was (and still is) scattered accross
several different source management systems. Few of them had the fancy built in
CI/CD features we all know and love. Further complicating matters, our build
system is old. Some parts of it are older than I am. SAS also delivers most of
its software rather than hosting it as a service. Not only do we have to support
the latest version, but any other supported versions we've shipped to customers.
To confront our myriad of problems, we needed an ecosystem agnostic solution
that was simple enough to work just about anywhere.

Our solution was Event Provenance Registry, or rather, the precursor to it. We
compared it to duct taping a Raspberry Pi to a rusty tractor. EPR is the glue
that enabled the rest of the pipeline to really take off.

## How it works

EPR is fairly simple in its operation, but requires some explanation. At a high
level, EPR collects events based on tasks done by the build pipeline and sends
them to a message queue. Other services that we call "watchers" watch the
message queue and take action when they see events of interest. EPR can gate
events by certain criteria as well. EPR supports Redpanda and Kafka as message
queues. For the examples below, I will assume we're using Redpanda.

To facilitate event collection, EPR has three data structures of note: events,
event-receivers, and event-receiver-groups. I will also refer to the latter two
as "receivers" and "groups" respectively for brevity.

### NVRPP

NVRPP is an unpronounceable acronym that you'll need to be familiar with to
truly understand how EPR works. It stands for:

- Name
- Version
- Release
- Package
- Platform ID

Each of these fields is just a string, though we strongly recommend you impose
some standards for how each is formatted, depending on your
situation. [Events](#events) requires these five fields to be defined. NVRPP is
based off of
the [NEVRA](https://docs.fedoraproject.org/en-US/modularity/core-concepts/nsvca/)
from Fedora. It allows us to represent most types of artifacts that might flow
through our pipeline. Events that have matching NVRPPs are associated with the
same artifact. This allows us to trace the flow of any artifact through our
pipeline, so long as events are posted at each step.

### Event Receivers

Event receivers are data structures stored within EPR that represent some kind
of action (i.e. a build, running a test, packaging an artifact, deploying a
binary, etc...). That have no dependencies and are classified by their name,
type, and version. You might name a receiver by the action it represents
like `golang-build-complete`. Types might be things like `build.finished`
or `artifact.packaged`. Receivers may have multiple [events](#events)
that correspond with them. Any events associated with a receiver must have a
payload that complies with the schema defined on the receiver. This allows some
guarantees about what kind of data you can expect of events going to any given
receiver.

```json
{
  "name": "golang-build-complete",
  "type": "build.finished",
  "version": "1.0.0",
  "description": "Receiver for Golang build results.",
  "enabled": true,
  "schema": {
    "type": "object",
    "properties": {
      "database": {
        "type": "string"
      }
    }
  }
}
```

When an event has been posted to a receiver, EPR will emit an event to Redpanda.

### Events

Events are a record of some action that took place in your pipeline and whether
it was successful. Each event is linked to a receiver by way of an ID. Events
are strictly formatted at the root level, with a free-form payload field that is
validated against the schema of its receiver. Each event contains a
boolean `success` field that represents whether an action was successful or not.
It also houses the [NVRPP](#nvrpp), which allows us to trace artifacts through
EPR.

When EPR receives an event, it posts the event and some receiver and group data
to Redpanda. Downstream watchers can then consume these messages and take their
own actions. Watchers can match messages based on the `success` field of a
message. This allows you to take different actions depending on if an event
passed or failed. For example, you could open a ticket against a team if their
event to the `artifact.scanned` type receiver had `success=false`.

```json
{
  "name": "my-app",
  "version": "2.2.2",
  "release": "2023-12-08:12-00-00",
  "platform_id": "linux",
  "package": "docker",
  "description": "This is an event for our application build.",
  "payload": {
    "database": "postgres"
  },
  "success": true,
  "event_receiver_id": "01HDS785T0V8KTSTDM9XGT33QQ"
}
```

### Event Receiver Groups

Event receiver groups can be thought of as gates that control whether an
artifact advances through the pipeline. Each group comprises multiple receivers.
Like receivers, groups can cause the generation of Redpanda events. However,
they only do this if each receiver has an event with a matching NVRPP
where `success=true`. Since there may be multiple events of varying successes
per receiver, only the most recent is considered. This allows you to run many
tasks in parallel, but only advance your artifact through the pipeline once all
its tasks have completed successfully.

```json
{
  "name": "release-checks",
  "type": "artifact.release",
  "version": "3.3.3",
  "description": "Send an event to release our application if all pipeline tasks have passed.",
  "enabled": true,
  "event_receiver_ids": [
    "01H9GW7FYY4XYE2R930YTFM7FM",
    "01HDS785T0V8KTSTDM9XGT33QQ"
  ]
}
```

- Explain watchers
- What it looks like in production

## Running EPR in Production

## Pitfalls

- No rbac for gates
- Adoption was difficult. Developers didn't like the box of legos approach.
  Challenge was as much political as technical.
- Laziness with gate schemas caused problems later.

## Benefits

- Greatly improved automated testing
- Automated software promotions
- Automated security scanning
- Allows groups to track the movement of artifacts through the pipline
