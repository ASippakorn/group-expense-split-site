# 07 - Show Deterministic Suggested Transfers

**What to build:** The Group detail view shows minimum-number Suggested Transfers using deterministic ordering, based on current Balances after Expenses, Splits, and Settlements.

**Blocked by:** 06 - Record and Delete Settlements.

**Status:** ready-for-agent

- [ ] Suggested Transfers are calculated from current Balances.
- [ ] Suggested Transfers minimize the number of repayments.
- [ ] Tie-breaking is deterministic so the same Balances always produce the same Suggested Transfers.
- [ ] Overpayment and zero-Balance cases are handled.
- [ ] The Group detail view displays Suggested Transfers clearly.
- [ ] OpenAPI documents the Suggested Transfer response shape.
- [ ] Backend service tests cover simple settlement, multiple debtors and creditors, zero Balances, overpayment, and deterministic ordering.
- [ ] Frontend tests cover rendering Suggested Transfers and the empty state.
