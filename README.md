# GoChat

GoChat is a simple realtime TCP chat application written in Go.  
It was built as a learning project to explore idiomatic Go patterns for concurrency, networking, and protocol design using only the standard library.

## Overview

The server accepts multiple TCP connections and coordinates communication through a central **hub**. The hub is responsible for registering and unregistering clients and broadcasting messages to all connected peers. Each client connection is handled using separate goroutines for reading from and writing to the network, ensuring non-blocking and isolated I/O.

The client connects to the server over TCP, captures user input from standard input, and sends messages encoded as newline-delimited JSON. Incoming messages from the server are processed concurrently and rendered asynchronously, keeping the client responsive while handling network events.

Both client and server share a small **protocol** package that defines the message format exchanged over the wire, avoiding duplication and keeping the communication contract explicit and consistent.

## Design Goals

- Use idiomatic Go concurrency with goroutines and channels  
- Keep a clear separation between networking, coordination, and rendering  
- Avoid shared mutable state in favor of message passing  
- Keep the protocol simple and explicit  

## Running the project

The project provides a simple `Makefile` for convenience:

```sh
make run-server
make run-client
```

Start the server first, then run one or more clients in separate terminals.

## Roadmap

This project is intentionally kept small, but several improvements are planned or considered as learning exercises:

* Commands support
Add client-side commands such as /quit, /nick, or /help.

* Chat rooms / channels
Allow users to join multiple rooms and broadcast messages selectively.

* Protocol evolution
Introduce message types (e.g. join, leave, system messages) and versioning.

* Persistence
Store chat history or user sessions using a simple storage layer.

* Improved client UX
Better terminal rendering, message formatting, and input handling.

* Graceful shutdown and reconnection
Handle disconnects and shutdown signals more robustly on both client and server.

## Notes

This project prioritizes clarity and correctness over feature completeness. It is intended as a foundation for experimentation and learning, and as a reference for building small concurrent systems in Go.
