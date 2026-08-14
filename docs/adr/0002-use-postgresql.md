# Use PostgreSQL for Persistence

We will use PostgreSQL as the system database because the project needs durable relational records for Users, Groups, Participants, Expenses, Splits, Settlements, sessions, and later currency rates. PostgreSQL also runs cleanly in Docker Compose and Kubernetes, which keeps the local and deployment story aligned for the course checkpoints.
