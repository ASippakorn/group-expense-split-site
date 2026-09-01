# 03 - Calculate Group Balances

**What to build:** Balances are derived from Expenses and Splits, using integer minor units and deterministic equal-Split rounding. Group cards and the Group detail view show Balance summaries.

**Blocked by:** 02 - Create Expense Ledger with Equal Splits.

**Status:** resolved

- [x] Balances are calculated from Expense and Split records rather than stored as source of truth.
- [x] Balance calculation handles each Participant's paid amount and owed amount.
- [x] Equal Split rounding produces stable Balance results.
- [x] Group cards show a lightweight Balance summary.
- [x] The Group detail view shows each Participant's Balance.
- [x] OpenAPI documents the Balance response shape.
- [x] Backend service tests cover one payer, multiple Participants, uneven rounding, and multiple Expenses.
- [x] Frontend tests cover Balance rendering on Group cards and Group detail.
