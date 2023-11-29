# EPR

Event Provenance Registry is culmination of several years of SAS's effort to
convert from large ship events to CI/CD. We built the first version internally
as a way to facilitate CI/CD in a complex, aging system. The result enables SAS
to build, package, scan, promote, and ship thousands of artifacts daily.

## Origin

Around mid 2019, management gave the directive that R&D needed to move faster to
stay competitive. SAS's software release process to that point consisted of two
or three large ship dates per year. Updates were slow and painful. To reverse
that trend, we needed a way to move artifacts through the pipeline more quickly.
The challenge was twofold. First, create a system that could allow disparate
pieces of our pipeline to communicate and chain together. Second, help R&D shift
gears from the old software development model to CI/CD.

- Why didn't we use fancy modern CI? Because we're a large company with large
  amounts of antiquated tooling. We needed something ecosystem agnostic.
- Allows groups to track the movement of artifacts through the pipline

## How it works

- Explain NVRPP
- Explain receivers, events, groups
- Explain watchers
- What it looks like in production

## Pitfalls

- No rbac for gates
- Adoption was difficult. Developers didn't like the box of legos approach.
  Challenge was as much political as technical.
- Laziness with gate schemas caused problems later.

## Benefits

- Greatly improved automated testing
- Automated software promotions
- Automated security scanning
