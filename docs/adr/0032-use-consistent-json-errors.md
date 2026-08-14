# Use Consistent JSON Errors

The API will return errors in a consistent JSON envelope with a machine-readable code, human-readable message, and optional field errors. This keeps frontend handling, tests, and OpenAPI documentation aligned across validation, authentication, authorization, and domain failures.
