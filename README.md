# Async Event Infrastructure Server

## Overview

A server for accepting events, storing events, and producing messages on a
message bus.

## Build

```bash
make
```

## Installation

To install in `/usr/local/bin` directory

Linux

```bash
make install
```

Mac OS X

```bash
make install-darwin
```

To install in your go path directory set `PREFIX`

Linux

```bash
make PREFIX=$(go env GOPATH) install
```

Mac OS X

```bash
make PREFIX=$(go env GOPATH) install-darwin
```

## Contributing

We welcome your contributions! Please read [CONTRIBUTING.md](CONTRIBUTING.md)
for details on how to submit contributions to this project.

## License

This project is licensed under the [Apache 2.0 License](LICENSE).
