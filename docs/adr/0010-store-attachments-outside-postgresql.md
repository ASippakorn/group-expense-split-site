# Store Attachments Outside PostgreSQL

Receipt Attachments and later voice notes will be stored in local filesystem or local object storage, with metadata and ownership records kept in PostgreSQL. This avoids bloating the relational database while still keeping the project fully local and portable through Docker Compose and Kubernetes volumes.
