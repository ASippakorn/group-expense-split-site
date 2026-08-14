# Use a Compose Test Database First

Integration tests will use a PostgreSQL test database reachable through Docker Compose before introducing Testcontainers. This keeps the first CI and local test setup simple while still exercising the same database engine used in development.
