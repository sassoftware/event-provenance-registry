# Use REPS for Enhancement Proposals

Many changes can be managed through Pull Requests.

However, some change requests are deemed "substantial," and these are subjected
to a design process to derive a consensus among the core maintainers.

A REP should include a GitHub Issue.

## Do I have to use the REP process?

Yes, if you want to see a large change adopted.

The REP process supports our new ways of working and creates a useful record of
enhancement requests that can help with change tracking.

## Template

```text
# Title

## Summary

### Motivation

### Goals

### Non-Goals

## Proposal

## Risks and Mitigation

## Graduation Criteria
```

## Prior Art

The REP process as proposed was essentially taken from the
[Kubernetes KEP Process](https://github.com/kubernetes/enhancements/tree/master/keps),
which was in turn adapted from the
[Rust RFC process](https://github.com/rust-lang/rfcs). The Rust process itself
seems to be very similar to the
[Python PEP process](https://www.python.org/dev/peps/pep-0001).
