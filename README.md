# Splitr

Splitr is a group expense splitter for trips, meals, and recurring gatherings. The first vertical slice includes registration, login, server-side sessions, Group creation, and Group listing.

## Stack

- Web: Vite, React, TypeScript, React Router, Tailwind CSS
- API: Go, Fiber, GORM
- Database: PostgreSQL with SQL migrations
- Tests: Go `testing` + Testify, Vitest, Playwright
- Delivery: Docker Compose, GitLab CI, Jenkinsfile, plain Kubernetes manifests

## Local Development

1. Copy `.env.example` to `.env`.
2. Start the local stack:

```powershell
docker compose up --build
```

3. Open the web app at `http://localhost:5173`.
4. API health is available at `http://localhost:8080/api/v1/health`.

## Repository Layout

- `apps/api` - Go Fiber API
- `apps/web` - React TypeScript app
- `docs/api` - OpenAPI contract
- `deploy` - Kubernetes manifests for Minikube
- `db/migrations` - SQL migrations
