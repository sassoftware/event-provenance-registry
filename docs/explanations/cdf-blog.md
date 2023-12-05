# EPR

Event Provenance Registry is culmination of several years of SAS's effort to
convert from large ship events to CI/CD. We built the first version internally
as a way to facilitate CI/CD in a complex, aging build system. The result enables SAS
to build, package, scan, promote, and ship thousands of artifacts daily.

## Origin

Around mid 2019, management gave the directive that R&D needed to move faster to
stay competitive. SAS's software release process to that point consisted of two
or three large ship events per year. Updates were slow and painful. To reverse
that trend, we needed a way to shorten the development cycle and deliver artifacts more quickly.
The challenge was twofold. First, create a system that could allow disparate
pieces of our pipeline to communicate and chain together. Second, help R&D shift
gears from the old software development model to CI/CD.

In a happy, imaginary world somewhere, we could have chained together Github actions or equivalent into a working pipeline. No need to write EPR at all. Reality was not so kind. Our source code was (and still is) scattered accross several different source management systems. Few of them had the fancy built in CI/CD features we all know and love. Further complicating matters, our build system is old. Some parts of it are older than I am. SAS also delivers most of its software rather than hosting it as a service. Not only do we have to support the latest version, but any other supported versions we've shipped to customers. To confront our myriad of problems, we needed an ecosystem agnostic solution that was simple enough to work just about anywhere.

Our solution was Event Provenance Registry, or rather, the precursor to it. We compared it to duct taping a Raspberry Pi to a rusty tractor. EPR is the glue that enabled the rest of the pipeline to really take off.

## How it works

EPR is fairly simple in its operation, but requires some explanation. At a high level, EPR collects events based on tasks done by the build pipeline and sends them to a message queue. Other services that we call "watchers" watch the message queue and take action when they see events of interest. EPR can gate events by certain criteria as well.

To facilitate event collection, EPR has three data structures of note: events, event-receivers, and event-receiver-groups. I will also refer to the latter two as "receivers" and "groups" respectively for brevity.

// TODO: need data structures here.
- event-receiver: An event-receiver is a 
- event:
- event-receiver-group:

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
- Allows groups to track the movement of artifacts through the pipline
