Status: ready-for-agent

# Group Expense Splitter Spec

## Problem Statement

Groups that share meals, trips, and recurring activities need a reliable way to record who paid, who participated, who owes whom, and which repayments have already happened. The project must be runnable locally, demonstrate realistic DevOps practices, and grow from the current registration and Group scaffold into a full Group Expense Splitter with tests, Docker Compose, CI/CD, Kubernetes deployment, and observability.

The current app has a Go/Fiber/GORM API, PostgreSQL schema for Users, sessions, Groups, and Participants, and a React/Vite web app for registration, login, Group creation, and Group listing. It does not yet support adding Participants by email, Expenses, Splits, Settlements, Balances, Suggested Transfers, or the Phase 2 features.

## Solution

Build Splitr as a monorepo application with a React TypeScript frontend, a Go Fiber backend, and PostgreSQL persistence. The product will let registered Users create Groups, add registered Users as Participants by email, record Expenses with equal, manual, percentage, and Tag-based Splits, calculate Balances from ledger records, record trusted Settlements, and show deterministic Suggested Transfers.

The system will treat Expenses, Splits, and Settlements as the source of truth. Balances and Suggested Transfers are calculated outputs, not manually maintained records. Money is stored as integer minor units using THB/satang by default, with original currency stored from the start so configurable Exchange Rates can be added in Phase 2.

The first buildable phase should complete the core financial loop after the existing auth/Group scaffold: manage Participants, create Expenses, calculate Balances, record Settlements, and cover the happy path with API, service, frontend, and Playwright tests.

## User Stories

1. As a User, I want to register with email and password, so that I can use Splitr without an external auth provider.
2. As a User, I want to sign in with email and password, so that I can access my Groups.
3. As a User, I want my session stored in an HTTP-only cookie, so that the frontend does not need to manage tokens.
4. As a User, I want to sign out, so that my session is removed.
5. As a User, I want to see only Groups where I am a Participant, so that private Group activity is not discoverable by others.
6. As a User, I want to create a Group with a name, Default Currency, and optional description, so that I can start tracking a shared trip, meal, or recurring gathering.
7. As an Owner, I want to add one registered User to a Group by email, so that they can participate in shared Expenses.
8. As an Owner, I want a clear "user not found" error when an email is not registered, so that I know why a Participant was not added.
9. As an Owner, I want to prevent adding the same User to a Group twice, so that the Participant list stays clean.
10. As an Owner, I want to mark a Participant inactive rather than deleting them when history exists, so that old Expenses, Splits, Settlements, and Balances remain explainable.
11. As a Regular Participant, I want to view the other active Participants in a Group, so that I know who can be included in new Expenses.
12. As a Regular Participant, I want inactive Participants shown in historical records, so that old Balances still make sense.
13. As a Participant, I want to add an Expense with amount, Payer, description, date, currency, and Participants, so that a shared cost is recorded.
14. As a Participant, I want the Expense Payer to be required as one of the Expense Participants, so that the Split cannot point outside the shared cost.
15. As a Participant, I want to create an equal Split, so that simple shared Expenses are fast to enter.
16. As a Participant, I want equal Split rounding to be deterministic, so that the same Expense always recalculates the same way.
17. As a Participant, I want to create a manual amount Split, so that I can handle uneven costs.
18. As a Participant, I want manual amount Splits rejected unless they exactly match the Expense amount after rounding, so that Balances remain correct.
19. As a Participant, I want to create a percentage Split, so that I can divide an Expense by agreed proportions.
20. As a Participant, I want percentage Splits rejected unless they total 100 percent, so that the allocation is complete.
21. As a Participant, I want to create Group-level Tags, so that Tag-based Splits can represent cases like "Alcohol" only charging selected Participants.
22. As a Participant, I want to select which Participants belong to a Tag for an Expense, so that the tagged Split charges only the eligible Participants.
23. As a Participant, I want Tags to be scoped to a Group, so that each Group can use its own sharing language.
24. As a Participant, I want to edit an Expense I created, so that data-entry mistakes can be fixed.
25. As an Owner, I want to edit or delete Expenses in my Group, so that the Group can correct mistakes.
26. As a Participant, I want deleted Expenses to be soft-deleted, so that financial history can be recovered and audited.
27. As a Participant, I want to see each Expense's Split details, so that I can understand why I owe or am owed money.
28. As a Participant, I want to see a Balance for each Participant, so that I know who is currently owed money and who currently owes money.
29. As a Participant, I want Balances calculated from Expenses, Splits, and Settlements, so that corrections automatically update the result.
30. As a Participant, I want to see raw net Balances, so that the calculation is transparent.
31. As a Participant, I want to see Suggested Transfers, so that the Group can settle with fewer repayments.
32. As a Participant, I want Suggested Transfers to minimize the number of transfers, so that settling is convenient.
33. As a Participant, I want Suggested Transfers to be deterministic, so that tests and explanations are stable.
34. As a Participant, I want either side of a repayment to record a Settlement, so that repayments can be tracked without an approval workflow.
35. As a Participant, I want a Settlement to record payer, receiver, amount, currency, date, and note, so that repayments are clear.
36. As a Participant, I want a Settlement amount to be allowed even when it exceeds the currently suggested amount, so that overpayment and real-world corrections are possible.
37. As a Participant, I want to delete and recreate a Settlement instead of editing it in place, so that repayment history stays simple.
38. As a Participant, I want Group cards to show lightweight Balance summaries, so that I can decide which Group needs attention.
39. As a Participant, I want the Group detail page to show Participants, Expenses, Balances, and Suggested Transfers, so that the main workflow is in one place.
40. As a User, I want default currency to be THB for Phase 1, so that local money entry is simple.
41. As a User, I want Expense currency stored from the start, so that Phase 2 Exchange Rates can be added without changing old records.
42. As a developer, I want all request validation enforced by the backend, so that invalid frontend behavior cannot corrupt data.
43. As a frontend user, I want client-side validation and accessible error messages, so that forms are easier to complete.
44. As a developer, I want API errors returned in a consistent JSON envelope, so that frontend handling and tests are predictable.
45. As a developer, I want the OpenAPI contract updated for each new route, so that the frontend, tests, and reviewers can see the API shape.
46. As a developer, I want CI to validate the OpenAPI contract, so that documentation drift is caught early.
47. As a teammate, I want Docker Compose to run PostgreSQL, the API, and the web app, so that local review takes one command.
48. As an instructor, I want the app to be runnable in a pipeline, so that checkpoints can verify real behavior.
49. As an instructor, I want unit, integration, and E2E tests, so that regressions are caught as requirements evolve.
50. As a developer, I want Demo Data for local development, so that I can review flows without repetitive setup.
51. As a developer, I want SQL migrations checked into the repo, so that schema changes are reviewable and repeatable.
52. As a developer, I want GORM models to match SQL migrations, so that persistence code and schema stay aligned.
53. As a developer, I want timestamps stored in UTC, so that reports and tests behave consistently across environments.
54. As a User, I want forms to be keyboard accessible, so that the app is usable without a mouse.
55. As a User, I want visible focus states and labeled inputs, so that I can navigate and correct forms confidently.
56. As a User, I want status cues that do not rely on color alone, so that important information is clear.
57. As a Phase 2 User, I want configurable Exchange Rates with CSV import, so that multi-currency Groups can be supported locally.
58. As a Phase 2 User, I want a Summary Export PDF with Group name, date range, Expenses, Settlements, Balances, and Suggested Transfers, so that I can share or archive the Group result.
59. As a Phase 2 User, I want a PromptPay Receiving ID on my profile, so that Suggested Transfers can show payment details.
60. As a Phase 2 User, I want Receipt Attachments stored with Expense metadata, so that supporting evidence can be reviewed.

## Implementation Decisions

- The app name is Splitr for the UI and `splitr` for code, package, and module naming.
- The backend uses Go, Fiber, GORM, PostgreSQL, and SQL migrations.
- The frontend uses Vite, React, TypeScript, React Router, Tailwind CSS, pnpm, and a small local component set.
- The repository remains a monorepo with separate web, API, deployment, database, API documentation, and project documentation areas.
- API routes are versioned under `/api/v1`.
- API behavior should remain REST-resource-oriented, with explicit domain action endpoints only where CRUD would be awkward.
- API errors use a consistent JSON envelope with code, message, and optional field errors.
- OpenAPI starts as a checked-in YAML contract and must be validated in CI. As routes stabilize, add tooling or generation to keep the contract synchronized with Fiber routes.
- Configuration comes from environment variables, with local defaults documented through `.env.example` and Docker Compose.
- User identity is unique email in Phase 1. Email change and password reset are out of scope for Phase 1.
- Passwords use Argon2id with per-user salts and a configurable pepper.
- Sessions are server-side records stored in PostgreSQL and sent through HTTP-only cookies.
- Sessions last seven days and are removed on logout.
- Users can see only Groups where they are Participants.
- Groups have an Owner and Regular Participants. More role levels are out of scope for Phase 1.
- Owners can add one registered User to a Group by email in Phase 1. Batch adding or import is deferred.
- Participants with historical records become Inactive Participants instead of being hard-deleted.
- Groups and Expenses are soft-deleted.
- Money is stored as integer minor units. THB/satang is the Phase 1 default.
- Every Expense stores original currency from the start.
- Balances are derived from Expenses, Splits, and Settlements rather than stored as source of truth.
- Expenses support equal, manual amount, manual percentage, and Tag-based Splits.
- Manual amount Splits must exactly match the Expense amount after rounding.
- Percentage Splits must total 100 percent.
- Tag-based Splits use Group-level Tags and charge only selected Participants.
- The Expense Payer must be included in the Expense Participants.
- Either side of a repayment can record a trusted Settlement in Phase 1.
- Settlements can be deleted and recreated, but not edited in place.
- Suggested Transfers minimize the number of repayments and use deterministic ordering.
- Receipt Attachments and voice notes are deferred to Phase 2 and should be stored outside PostgreSQL, with metadata in PostgreSQL.
- PromptPay Receiving IDs belong to Users, not Groups.
- Summary Export is a Phase 2 PDF report over a date range, including Expenses, Settlements, Balances, and Suggested Transfers.
- Observability in Phase 2 uses Prometheus, Grafana, Loki, Promtail, and structured JSON logs from the API.
- Docker Compose should run PostgreSQL, the API, and the web app.
- Kubernetes deployment starts with plain manifests for Minikube.
- GitLab CI is the primary pipeline-as-code path, with Jenkinsfile support for course verification.

## Testing Decisions

- Tests should lock down external behavior and business rules, not implementation details.
- The primary backend seam is service-level tests around Participants, Expenses, Splits, Settlements, Balances, and Suggested Transfers.
- Money tests must cover integer minor units, equal Split rounding, manual amount validation, percentage validation, Tag-based Split eligibility, Settlement application, overpayment, and deterministic Suggested Transfers.
- API tests should cover authenticated and unauthenticated behavior, Group visibility, Participant add-by-email, Expense creation, Settlement recording, soft deletion behavior, and consistent JSON errors.
- Frontend tests should cover page/form behavior, accessible labels, validation messages, loading states, and error states.
- Playwright should cover the happy-path user story: register, create Group, add Participant, add Expense, view Balance, record Settlement, and observe updated Balance.
- OpenAPI validation should run in CI and fail when the required route contract is missing or malformed.
- Existing prior art includes the API password hash unit test, the web app sign-in render test, and the first OpenAPI validation script.
- Integration tests should use a PostgreSQL test database through Docker Compose first. Testcontainers can be introduced later only if useful.

## Out of Scope

- External auth providers.
- External databases, cloud storage, or SaaS APIs.
- Password reset in Phase 1.
- Email change in Phase 1.
- Guest or non-login Participants in Phase 1.
- Batch Participant email import in Phase 1.
- Receipt Attachment upload in Phase 1.
- Voice notes in Phase 1.
- Live currency conversion APIs.
- PromptPay QR generation and payslip recording in Phase 1.
- Summary Export PDF in Phase 1.
- Offline mode and sync.
- OCR receipt parsing.
- Additional Group roles beyond Owner and Regular Participant.
- Full browser support beyond current Chrome and Edge for coursework review.
- Helm charts unless plain Kubernetes manifests become painful.

## Further Notes

- Use the glossary in `CONTEXT.md` consistently: User, Group, Participant, Owner, Regular Participant, Inactive Participant, Expense, Split, Tag, Settlement, Balance, Suggested Transfer, Exchange Rate, Receipt Attachment, Summary Export, PromptPay Receiving ID, and Demo Data.
- The first implementation tickets should be tracer bullets that keep the app runnable after each slice.
- The next best sequence is: Participant add-by-email, Expense/Split/Balance backend, Expense UI, manual and percentage Split validation, Tag-based Split, Settlement recording, Suggested Transfers, full Docker Compose verification, and Playwright E2E coverage.
- The current folder is not a Git repository yet. Initialize Git before relying on diff-based review or commit workflows.
