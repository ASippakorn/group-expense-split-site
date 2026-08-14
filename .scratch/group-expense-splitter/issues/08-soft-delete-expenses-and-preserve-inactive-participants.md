# 08 - Soft Delete Expenses and Preserve Inactive Participants

**What to build:** Owners can soft-delete Expenses and mark historical Participants inactive without breaking old financial records.

**Blocked by:** 06 - Record and Delete Settlements.

**Status:** ready-for-agent

- [ ] Owners can soft-delete Expenses in their Group.
- [ ] Deleted Expenses no longer affect current Balances.
- [ ] Deleted Expenses remain recoverable for audit/debug purposes.
- [ ] Owners can mark a Participant inactive when the Participant has historical records.
- [ ] Inactive Participants remain visible in historical Expenses, Splits, Settlements, and Balances where needed.
- [ ] Inactive Participants cannot be selected for new Expenses.
- [ ] OpenAPI documents Expense deletion and Participant deactivation behavior.
- [ ] Backend tests cover soft deletion, Balance recalculation, inactive Participant visibility, and new Expense exclusion.
- [ ] Frontend tests cover deleting an Expense and deactivating a Participant.
