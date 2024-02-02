# Use ADRs

## History

- Status: accepted
- Decision makers: Brett Smith
- Date: 20231103

The following text is almost directly replicated from a
[Michael Nygard blog post](http://thinkrelevance.com/blog/2011/11/15/documenting-architecture-decisions).
This allows us to reference it more directly without always going to a website
to refresh.

## Context and problem statement

Architecture for agile projects has to be described and defined differently. Not
all decisions are made at once, nor are all of them complete when the project
begins.

Agile methods are not opposed to documentation, only to valueless documentation.
Documents that assist the team itself can have value, but only if they are kept
up to date. Large documents are never kept up to date. Small, modular documents
have at least a chance at being updated.

Nobody ever reads large documents, either. Most developers have been on at least
one project where the specification document was larger (in bytes) than the
total source code size. Those documents are too large to open, read, or update.
Bite sized pieces are easier for for all stakeholders to consume.

One of the hardest things to track during the life of a project is the
motivation behind certain decisions. A new person coming on to a project may be
perplexed, baffled, delighted, or infuriated by some past decision. Without
understanding the rationale or consequences, this person has only two choices:

1. Blindly accept the decision.

   This response may be OK, if the decision is still valid. It may not be good,
   however, if the context has changed and the decision should really be
   revisited. If the project accumulates too many decisions accepted without
   understanding, then the development team becomes afraid to change anything
   and the project collapses under its own weight.

2. Blindly change it.

   Again, this may be OK if the decision needs to be reversed. On the other
   hand, changing the decision without understanding its motivation or
   consequences could mean damaging the project's overall value without
   realizing it. (e.g., the decision supported a non-functional requirement that
   hasn't been tested yet.)

It's better to avoid either blind acceptance or blind reversal.

## Decision outcome

We will keep a collection of records for "architecturally significant"
decisions: those that affect the structure, non-functional characteristics,
dependencies, interfaces, or construction techniques.

An architecture decision record is a short text file in a format similar to an
Alexandrian pattern. (Though the decisions themselves are not necessarily
patterns, they share the characteristic balancing of forces.) Each record
describes a set of forces and a single decision in response to those forces.
Note that the decision is the central piece here, so specific forces may appear
in multiple ADRs.

We keep ADRs in the project repository under doc/arch/adr-NNN.md

We use a lightweight text formatting language like Markdown or Textile.

ADRs are numbered sequentially and monotonically. Numbers are not reused.

If a decision is reversed, we keep the old ADR, but mark it as superseded. (It's
still relevant to know that it was the decision, but is no longer the decision.)

A template format for these decisions can be found in the doc/arch directory.

### Consequences

One ADR describes one significant decision for a specific project. The decision
should be something that has an effect on how the rest of the project runs.

The consequences of one ADR are very likely to become the context for subsequent
ADRs. This is also similar to Alexander's idea of a pattern language: the
large-scale responses create spaces for the smaller scale to fit into.

Developers and project stakeholders can see the ADRs, even as the team
composition changes over time.

The motivation behind previous decisions is visible for everyone, present and
future. Nobody is left scratching their heads to understand, "What were they
thinking?" and the time to change old decisions will be clear from changes in
the project's context.

## Links <!-- optional -->

None
