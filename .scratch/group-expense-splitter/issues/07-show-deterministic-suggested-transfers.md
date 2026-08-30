# 07 - Show Deterministic Suggested Transfers

**What to build:** The Group detail view shows minimum-number Suggested Transfers using deterministic ordering, based on current Balances after Expenses, Splits, and Settlements.

**Blocked by:** 06 - Record and Delete Settlements.

**Status:** claimed

- [x] Suggested Transfers are calculated from current Balances.
- [x] Suggested Transfers minimize the number of repayments.
- [x] Tie-breaking is deterministic so the same Balances always produce the same Suggested Transfers.
- [x] Overpayment and zero-Balance cases are handled.
- [x] The Group detail view displays Suggested Transfers clearly.
- [x] OpenAPI documents the Suggested Transfer response shape.
- [x] Backend service tests cover simple settlement, multiple debtors and creditors, zero Balances, overpayment, and deterministic ordering.
- [x] Frontend tests cover rendering Suggested Transfers and the empty state.
