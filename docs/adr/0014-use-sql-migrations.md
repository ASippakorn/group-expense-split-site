# Use SQL Migrations

Database changes will be managed with checked-in SQL migrations run through golang-migrate. This keeps PostgreSQL schema changes explicit, reviewable, and reproducible across local development, Docker Compose, CI, and Kubernetes deployments.
