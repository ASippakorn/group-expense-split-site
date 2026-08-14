# 09 - Expand Demo Data and Docker Compose Verification

**What to build:** Local Demo Data covers Users, a Group, Participants, Expenses, Splits, Settlements, Balances, and Suggested Transfers. Docker Compose can run the full stack for local and instructor review.

**Blocked by:** 07 - Show Deterministic Suggested Transfers; 08 - Soft Delete Expenses and Preserve Inactive Participants.

**Status:** ready-for-agent

- [ ] Demo Data includes multiple Users, at least one Group, active Participants, Expenses, Splits, Settlements, Balances, and Suggested Transfers.
- [ ] Demo Data can be loaded repeatably without duplicating records.
- [ ] Docker Compose starts PostgreSQL, the API, and the web app.
- [ ] Migrations run before the API starts.
- [ ] The web app can call the API in the Compose environment.
- [ ] README instructions explain how to run and verify the full stack locally.
- [ ] Verification includes at least one command or documented check for API health and web availability.
- [ ] CI configuration remains aligned with the local stack expectations.
