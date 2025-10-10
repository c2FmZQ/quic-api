[![Go](https://github.com/c2FmZQ/quic-api/actions/workflows/go.yml/badge.svg)](https://github.com/c2FmZQ/quic-api/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/c2FmZQ/quic-api.svg)](https://pkg.go.dev/github.com/c2FmZQ/quic-api)

# quic-api

This repo contains auto-generated interfaces and wrappers for the [quic-go](https://pkg.go.dev/github.com/quic-go/quic-go) data structures.

See https://github.com/c2FmZQ/tlsproxy/issues/211 for context and motivation.

## Why use interfaces for this?

Using interfaces in this context provides a few key benefits:

1.  **Decoupling from a Third-Party Library:** This repository helps decouple projects from the `quic-go` library. By depending on the stable interfaces in `quic-api` instead of the concrete types from `quic-go`, downstream projects can be insulated from breaking API changes in the underlying library.

2.  **Enabling Testability:** The generated interfaces allow for the creation of mock implementations of the `quic-go` data structures for unit testing. This allows for testing business logic without needing a live QUIC connection.

3.  **Following Go Idioms:** This approach allows consuming projects to follow the Go proverb: "Accept interfaces, return structs." Code can accept an interface from this library, rather than a concrete type from `quic-go`.

While a library provider would ideally provide stable interfaces, this repository offers a practical alternative when they are not available.
