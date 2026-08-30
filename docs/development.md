# Development guide

How to run Splitr locally, run the test suites, and pick up the next ticket.

## Prerequisites

- Docker (Docker Compose v2)
- For working on the API without Docker: Go 1.26+
- For working on the web app without Docker: Node 24+, pnpm 11 (`corepack enable`)

## Run the full stack

```powershell
docker compose up --build
```

- Web: http://localhost:5173
- API health: http://localhost:8080/api/v1/health
- Postgres: `localhost:5432` (user/password/database all `splitr`)

Migrations in `db/migrations` run automatically via the `migrate` service before the API starts. Postgres data persists in the `postgres-data` volume; use `docker compose down -v` to reset it.

## Testing the API manually

Auth is cookie-session based, so keep a cookie jar. From the repository root:

```bash
BASE=http://localhost:8080/api/v1
JAR=cookies.txt

curl -c $JAR -X POST $BASE/auth/register -H "Content-Type: application/json" \
  -d '{"email":"alice@test.com","password":"password123"}'

curl -b $JAR $BASE/me

curl -b $JAR -X POST $BASE/groups -H "Content-Type: application/json" \
  -d '{"name":"Phuket Trip","defaultCurrency":"THB","description":"August trip"}'

curl -b $JAR -X POST $BASE/groups/<groupId>/participants -H "Content-Type: application/json" \
  -d '{"email":"bob@test.com"}'

curl -b $JAR $BASE/groups/<groupId>/participants

# amountMinor is integer minor units: 150000 = THB 1500.00
curl -b $JAR -X POST $BASE/groups/<groupId>/expenses -H "Content-Type: application/json" \
  -d '{"description":"Dinner","amountMinor":150000,"currency":"THB","expenseDate":"2026-08-15",
       "payerParticipantId":"<payerParticipantUuid>","participantIds":["<uuid1>","<uuid2>"]}'

curl -b $JAR $BASE/groups/<groupId>/expenses
```

Participant IDs come from the participants listing. The full contract is in `docs/api/openapi.yaml`.

## Tests

```bash
# API unit tests (repository is faked; no database needed)
cd apps/api && go test ./...

# Web unit tests (Vitest)
cd apps/web && pnpm test

# Type-check
cd apps/web && pnpm lint

# E2E (Playwright; needs the stack running)
cd apps/web && pnpm exec playwright test

# OpenAPI contract lint
pnpm openapi:lint
```

## Conventions

- **Layering**: handlers (`internal/http`) → services (`internal/service`) → repository (`internal/repository`). Business rules and authorization belong in services; handlers only translate errors to HTTP. See `docs/architecture.md`.
- **Money**: always integer minor units (`amountMinor`), never floats.
- **Domain language**: use the terms defined in `CONTEXT.md` (Participant, not member; Expense, not bill).
- **API changes**: update `docs/api/openapi.yaml` in the same change and run `pnpm openapi:lint`.

## Ticket workflow

Tickets live as markdown files in `.scratch/group-expense-splitter/issues/` and are mirrored to GitHub Issues (see `docs/agents/issue-tracker.md`). Each ticket states what to build, its blocking dependencies, and its acceptance checklist.

1. Check statuses in the issue files; a ticket is ready when everything it is blocked by is `resolved`.
2. Create a branch named after the ticket (e.g. `split03-calculate-group-balances`).
3. Implement, keeping the acceptance checklist as the definition of done.
4. Run the full test suites and `pnpm openapi:lint`.
5. Commit and merge to `main`, then flip the ticket's status to `resolved`.

## Docker build notes

- The web image runs `pnpm install --frozen-lockfile`, so `pnpm-lock.yaml` must be committed whenever workspace manifests change.
- `node_modules` and `apps/web/dist` are excluded from the build context via `.dockerignore`; never copy host `node_modules` into an image (Windows junction symlinks break Linux builds).

## Current status

Implemented: tickets 01 (Participant management) and 02 (Expense ledger with equal splits). The next unblocked ticket is 03 (Calculate group balances); the full dependency chain is recorded in the ticket files.
