# Server-Sent Events (SSE) Package

This package provides functionality for working with Server-Sent Events (SSE) in ogen-generated clients.

## Overview

Server-Sent Events is a technology that allows a server to push real-time updates to a client over a single HTTP connection. This package implements:

1. `Reader` for parsing SSE streams
2. `Event` data structure for representing individual events

## Usage

When the OpenAPI specification includes a response with media type `text/event-stream`, ogen will generate client code that returns a channel of the response object defined in the schema.

### Example

For an endpoint with a response schema:

```yaml
responses:
  200:
    content:
      text/event-stream:
        schema:
          type: object
          properties:
            message:
              type: string
            timestamp:
              type: integer
```

The generated client will:

1. Return a channel of the response object
2. Handle the streaming nature of SSE
3. Automatically close the channel when the context is canceled or the connection is closed

## Implementation Details

- The SSE parsing follows the [W3C EventSource API standard](https://html.spec.whatwg.org/multipage/server-sent-events.html)
- Events are expected to contain JSON data (which will be parsed according to the schema)
- The client context can be used to cancel the stream
