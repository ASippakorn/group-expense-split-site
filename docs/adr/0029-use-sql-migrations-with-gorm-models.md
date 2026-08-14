# Use SQL Migrations with GORM Models

The API will use checked-in SQL migrations for schema changes while using GORM models for application persistence. This avoids relying on GORM AutoMigrate for repeatable environments and keeps database changes explicit for local development, CI, Docker Compose, and Kubernetes deployments.
