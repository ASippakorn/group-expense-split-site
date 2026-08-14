# 11 - Keep OpenAPI and CI Aligned

**What to build:** OpenAPI documents the implemented routes and CI validates the contract, backend tests, web checks, and E2E path.

**Blocked by:** 01 - Add Participant Management; 02 - Create Expense Ledger with Equal Splits; 04 - Support Manual Amount and Percentage Splits; 05 - Support Group Tags and Tag-Based Splits; 06 - Record and Delete Settlements; 07 - Show Deterministic Suggested Transfers; 08 - Soft Delete Expenses and Preserve Inactive Participants.

**Status:** ready-for-agent

- [ ] OpenAPI documents all implemented auth, Group, Participant, Expense, Split, Tag, Settlement, Balance, and Suggested Transfer routes.
- [ ] OpenAPI includes consistent error responses for validation, authentication, authorization, and not found cases.
- [ ] CI validates the OpenAPI contract.
- [ ] CI runs backend tests.
- [ ] CI runs web TypeScript checks.
- [ ] CI runs web unit tests.
- [ ] CI includes the E2E stage or documents the services required for it.
- [ ] The contract validation catches missing first-slice and financial workflow routes.
