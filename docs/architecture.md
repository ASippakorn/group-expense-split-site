# Architecture

This document explains how the Splitr codebase is organized and how a request flows through it. For domain terminology (User, Group, Participant, Expense, Split, Settlement), see `CONTEXT.md` at the repository root. For the reasoning behind individual technology choices, see the ADRs in `docs/adr/`.

## Stack

| Layer | Technology |
|---|---|
| Web | Vite, React, TypeScript, React Router, Tailwind CSS |
| API | Go, Fiber, GORM |
| Database | PostgreSQL with SQL migrations (`golang-migrate`) |
| Tests | Go `testing` + Testify, Vitest, Playwright |
| Delivery | Docker Compose locally, plain Kubernetes manifests for Minikube |

## Repository layout

```
apps/api/cmd/api/           entrypoint: wiring only
apps/api/internal/config/   environment-based configuration
apps/api/internal/database/ Postgres connection + migrate
apps/api/internal/domain/   domain structs and constants, no dependencies
apps/api/internal/http/     Fiber handlers, middleware, error responses
apps/api/internal/repository/  GORM data access (Store)
apps/api/internal/security/ password hashing
apps/api/internal/service/  business rules and authorization
apps/web/src/               React app (single-page, session cookie auth)
db/migrations/              versioned SQL migrations
docs/api/                   OpenAPI contract
docs/adr/                   architecture decision records
.scratch/<feature>/issues/  local issue tracker tickets
```

## API layering

The API follows a strict handler → service → repository layering. Each layer has one job:

- **Handlers** (`apps/api/internal/http/`) parse and validate the shape of a request (UUIDs, JSON body) and translate service errors into HTTP responses. They contain no business logic.
- **Services** (`apps/api/internal/service/`) own all business rules and authorization. They communicate failure through sentinel errors (`service.ErrForbidden`, `service.ErrValidation`, `service.ErrUserNotFound`, `service.ErrParticipantExists`), never HTTP codes.
- **Repository** (`apps/api/internal/repository/store.go`) owns all GORM/SQL. It returns `repository.ErrNotFound` and raw errors; it knows nothing about authorization.

A service depends on a small interface (`GroupRepository` in `service/groups.go`), which is what the service tests fake out. Domain structs live in `domain/models.go` and are shared across layers.

Routes are registered in `apps/api/internal/http/server.go`. Everything except `/health` and the auth routes sits behind the `requireUser` middleware.

## Authentication

Sessions are server-side (ADR 0005). The flow:

1. `POST /auth/register` or `/auth/login` verifies credentials (bcrypt, `internal/security/password.go`) and calls `store.CreateSession`, which inserts a row with a UUID token and TTL.
2. The session token is returned as a cookie. The browser sends it on every request; `apps/web/src/api.ts` uses `credentials: "include"`.
3. `requireUser` (`internal/http/middleware.go`) parses the cookie, loads the user via `store.FindSessionUser`, and rejects with 401 when the session is missing or expired.
4. `POST /auth/logout` deletes the session row.

## Request walkthrough: Owner adds a Participant

The rule from ticket 01: only the Owner can add Participants, and only registered Users can be added.

1. The browser form in `apps/web/src/App.tsx` (rendered only when the current user owns the group) calls `addParticipant()` in `apps/web/src/api.ts`, which sends `POST /api/v1/groups/{groupId}/participants` with `{"email": "..."}`.
2. `requireUser` resolves the session to a User ID.
3. `addParticipant` handler (`internal/http/group_handlers.go`) parses the group UUID and JSON body, then calls the service.
4. `GroupService.AddParticipantByEmail` (`internal/service/groups.go`) runs the rule chain:
   - `FindParticipant(groupID, actorID)` — the caller must be an active Owner of the group, otherwise `ErrForbidden`.
   - `normalizeEmail` — otherwise `ErrValidation`.
   - `FindUserByEmail` — the target must already be registered, otherwise `ErrUserNotFound`. There is no invite flow.
   - `FindParticipant(groupID, targetUserID)` — duplicates return `ErrParticipantExists`.
   - `CreateParticipant` inserts the row with role `participant` and `active = true`.
5. The handler maps the outcome to 201 with the embedded participant, or to 403 / 400 / 409 with an error code.
6. The `UNIQUE(group_id, user_id)` constraint on the `participants` table is the final defense against concurrent duplicate adds.

The Owner's own participant row is not created through this endpoint. `store.CreateGroupWithOwner` (`internal/repository/store.go`) creates the group and its Owner participant together at group creation, so every group has exactly one Owner from birth.

## Request walkthrough: Participant creates an Expense

The rules from ticket 02: any active Participant can create an Expense; the Payer must be one of the selected Participants; splits are equal; rounding is deterministic.

1. `GroupService.CreateEqualExpense` (`internal/service/expenses.go`) first checks the actor is an active Participant of the group (`requireActiveParticipant`).
2. Field validation: non-empty description, positive `amountMinor`, 3-letter currency (default `THB`), `YYYY-MM-DD` date.
3. Every selected participant ID must resolve to an active participant **of this group** (`FindParticipantByID`), otherwise `ErrValidation` — you cannot split with someone outside the group.
4. The Payer must be among the selected participants.
5. `equalSplits` divides `amountMinor` by the participant count. The base amount is `amountMinor / count`; the first `amountMinor % count` participants (in ascending participant-UUID order — the sort happens before splitting) each receive one extra minor unit. Splits therefore always sum exactly to the expense amount and the result is reproducible.
6. `store.CreateExpenseWithSplits` writes the expense and its splits.

## Money

Amounts are stored as integer minor units (`amountMinor`, e.g. `150000` = THB 1500.00) per ADR 0006. The API never accepts floating-point money. The frontend converts user input to minor units (`parseMoneyMinor` in `App.tsx`) and formats it back for display (`formatMoney`).

## Frontend

`apps/web/src` is a single-page React app with three routes: `/login`, `/register`, and the dashboard `/`. The dashboard contains group creation, the group list, and — when a group is selected — the Participants panel and the Expenses panel. All state is local to `App.tsx` via hooks; there is no global state library. All server calls go through `apps/web/src/api.ts`, which owns the base URL, credentials, and error extraction.

## What is implemented

Tickets 01 (Participant management) and 02 (Expense ledger with equal splits) are complete. Balances, settlements, non-equal split types, tags, and suggested transfers do not exist yet; see the ticket statuses in `.scratch/group-expense-splitter/issues/`. The dashboard group cards currently display placeholder metrics (`THB 0.00`) until balances land with ticket 03.
