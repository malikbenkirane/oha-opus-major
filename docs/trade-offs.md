# Trade-offs due to time limitation

## Inadequate Testing Practices

Testing is currently done manually. We'll evaluate if we can finish it and
identify priorities to meet the suggested timeline.

## GET /update-player-data endpoint implementation

The current implementation lacks support for CORS and cache headers; we would
incorporate those features if we had additional time.

The handler wrapper doesn’t handle errors that arise after part of the response
body has already been sent. We could have resolved this issue if we had more
time.

Storing the entire response in a list isn’t practical for large multiplayer
games because it can consume excessive memory. Ideally, the results should be
paginated, but we opted for this simpler solution due to time constraints.

## Error handling

For this assignment, we adopted a minimalist approach to error handling.

Errors are simply printed to the console, and no attempt is made to recover
from them or provide structured handling.

Additionally, server‑side errors are relayed directly to the client without any
filtering. While this speeds development, it is not suitable for production, as
exposing raw server errors can create security vulnerabilities and hinder
robust fault tolerance.

## Adapters options validation

**Current limitation**

The player‑data server adapter doesn’t validate the options it receives.
Consequently, a misconfigured parameter set can slip through to the port
initializer unnoticed, and the problem isn’t discovered until the server is
actually instantiated (e.g. an invalid server address).

**Proposed improvement**  

Validate the configuration parameters before constructing the server. By doing
so, mis‑specified options (e.g., invalid port numbers, missing required fields,
unsupported values) can be reported early, preventing runtime failures.

**API change**  
Update the signature of the `New` constructor to return an additional `error`
value alongside the instantiated port. For example:

    // Before
    func New(opts... Options) (port.Port, error)

    // After
    func New(opts... Options) (port.Port, error) // where the second return is the error

This change makes it explicit to callers that server creation can fail due to
invalid configuration, enabling them to handle the error immediately.
