# Use GORM for Database Access

The Go API will use GORM for database access because the project owner prefers an ORM-based workflow for building the app quickly. This trades some SQL explicitness for faster model-driven development, so financial queries and Balance calculations should still be covered with focused tests.
