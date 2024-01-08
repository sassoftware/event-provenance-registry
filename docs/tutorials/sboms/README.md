# CDEvents

## Overview

[OWASP CycloneDX](https://cyclonedx.org/specification/overview/) is a full-stack
Bill of Materials (BOM) standard that provides advanced supply chain
capabilities for cyber risk reduction.

In this tutorial we will learn how we can use CycloneDX SBOMs with the Event
Provenance Registry (EPR).

## Requirements

The [Hello World](../hello_world/README.md) has been completed and the EPR
server is running.

- [jq](https://jqlang.github.io/jq/)
- [cdxgen](https://github.com/CycloneDX/cdxgen)
- [syft](https://github.com/anchore/syft)
- [grype](https://github.com/anchore/grype)

## CycloneDX bom schema

Download the CycloneDX bom schema.

```bash
curl -ssLO https://raw.githubusercontent.com/CycloneDX/specification/1.5/schema/bom-1.5.schema.json
```

## Create a source SBOM event

First we will create the event receiver and apply the CycloneDX v1.5 schema for
artifact sbom created.

To make this easier we will create the JSON data we need to post first.

```bash
echo "{\"name\": \"artifact-cyclonedx-sbom\",\"type\": \"build.artifact.cyclonedx.sbom\",\"version\": \"1.0.0\",\"description\": \"Artifact CycloneDX v1.5 SBOMs\",\"enabled\": true,\"schema\": $(cat bom-1.5.schema.json)}" | jq > er.json
```

Create the event receiver:

```bash
curl --location --request POST 'http://localhost:8042/api/v1/receivers' \
--header 'Content-Type: application/json' \
--data @er.json
```

The results of the command should look like this:

```json
{ "data": "01HGBDPKPXVYKFNJ1Q6NDK7AMK" }
```

Create an SBOM to post.

Run the following command at the root of the event-providence-registry checkout.

```bash
cdxgen -o ./docs/tutorials/workshops/sboms/sbom.json --spec-version 1.5
```

Now we create the data for our event.

```bash
cd ./docs/tutorials/workshops/sboms
```

```bash
echo "{\"name\": \"epr\",\"version\": \"1.0.1\",\"release\": \"2023.11.16\",\"platform_id\": \"aarch64-gnu-linux-7\",\"package\": \"oci\",\"description\": \"scan source code for OCI image EPR\",\"payload\": $(cat sbom.json),\"success\": true,\"event_receiver_id\": \"01HGBDPKPXVYKFNJ1Q6NDK7AMK\"}" | jq > sbom_event.json
```

Now that we have the data ready we will POST the event to the event receiver.
The event payload will be in the form of a CycloneDX SBOM.

Create an event:

```bash
curl --location --request POST 'http://localhost:8042/api/v1/events' \
--header 'Content-Type: application/json' \
--data @sbom_event.json
```

The results of the command should look like this:

```json
{ "data": "01HGBDRPZRT96R4586GFA13W91" }
```

## Create an SBOM for an OCI image

Create an OCI image to scan

Run the following command at the root of the event-providence-registry checkout.

```bash
make docker-image
```

This command will create an OCI image for EPR called `epr-server:local`

We can create the SBOM for the image using `syft` as follows:

```bash
syft epr-server:local  --scope all-layers -o cyclonedx-json=./docs/tutorials/workshops/sboms/oci_sbom.json
```

Now we create the data for our event.

```bash
cd ./docs/tutorials/workshops/sboms
```

```bash
echo "{\"name\": \"epr\",\"version\": \"1.0.1\",\"release\": \"2023.11.16\",\"platform_id\": \"aarch64-gnu-linux-7\",\"package\": \"oci\",\"description\": \"scan of the EPR OCI image\",\"payload\": $(cat oci_sbom.json),\"success\": true,\"event_receiver_id\": \"01HGBDPKPXVYKFNJ1Q6NDK7AMK\"}" | jq > oci_sbom_event.json
```

Now that we have the data ready we will POST the event to the event receiver.
The event payload will be in the form of a CycloneDX SBOM.

Create an event:

```bash
curl --location --request POST 'http://localhost:8042/api/v1/events' \
--header 'Content-Type: application/json' \
--data @oci_sbom_event.json
```

The results of the command should look like this:

```json
{ "data": "01HGBDVCEWE5KYSNMYJPECQEYN" }
```

Now we can retrieve the SBOM and use a tool like `grype` to scan it for
vulnerabilities.

```bash
curl --location --request GET 'http://localhost:8042/api/v1/events/01HGBDVCEWE5KYSNMYJPECQEYN' \
--header 'Content-Type: application/json'  | jq .data.payload | grype
```
