# Event Provenance Registry Documentation

## Overview

The Event Provenance Registry (EPR) is a service that manages and stores events
and tracks event-receivers and event-receiver-groups. EPR provides an API that
lets you create events, event-receivers, and event-receiver-groups. You can
query the EPR using the GraphQL endpoint to get identifying information about
events, event-receivers, and event-receiver-groups. EPR collects events from the
supply chain to record the lifecycle of a unit in the SDLC. EPR validates
event-receivers have events. EPR emits a message when a event is received as
well as when an event-receiver-groups is complete for a unit version.

## Documentation

The Event Provenance Registry (EPR) documentation will follow the
[Diátaxis Framework](https://diataxis.fr) wherever it makes sense.

## Table of Contents

| Name                           | Description                                          |
| ------------------------------ | ---------------------------------------------------- |
| [Tutorials](./tutorials)       | Tutorials (getting started)                          |
| [How-to](./how-to)             | How-to Guides (step by step solutions)               |
| [Reference](./reference)       | Reference Documentation (descriptions of the system) |
| [Explanations](./explanations) | Explanations (ADR, Enhancement Requests, SIG)        |

## Sections

### Tutorials

Tutorials are learning oriented. The idea is to get the users started with
practical steps.

### How-tos

The how-to guides are problem/solution oriented. They take the reader through a
series of steps to solve a specific problem or create a solution.

### References

Reference documentation is information oriented. They represent the technical
descriptions of the machinery and how to run it. How to use Key classes,
functions, and APIs. Reference material explains things like "this is how the
GraphQL front end works".

### Explanations

Explanations are understanding oriented. They discuss the background and context
of the architecture, decisions made, research that was done, and the pros and
cons of things. Discussions explain why things are the way they are.
