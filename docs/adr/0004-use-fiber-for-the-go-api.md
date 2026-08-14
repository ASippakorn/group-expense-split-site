# Use Fiber for the Go API

We will use Fiber for the Go HTTP API because the project owner prefers it and it gives the backend a compact routing and middleware model for authentication, validation, and JSON APIs. This choice couples handlers to Fiber's request context, so future contributors should preserve that style unless there is a clear reason to migrate.
