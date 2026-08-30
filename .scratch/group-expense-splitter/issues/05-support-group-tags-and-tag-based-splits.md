# 05 - Support Group Tags and Tag-Based Splits

**What to build:** Participants can create Group-level Tags and use them to charge only selected Participants on tagged Expenses.

**Blocked by:** 04 - Support Manual Amount and Percentage Splits.

**Status:** resolved

- [x] Participants can create Tags scoped to a Group.
- [x] Participants can view the Tags for a Group.
- [x] A tagged Expense charges only the selected Participants for that Tag.
- [x] Tag-based Splits produce correct Balances.
- [x] Tags from one Group cannot be used in another Group.
- [x] The Group detail view supports creating and selecting Tags for an Expense.
- [x] OpenAPI documents Tag routes and Tag-based Split inputs.
- [x] Backend tests cover Tag scope, Tag membership, and Balance effects.
- [x] Frontend tests cover Tag creation and Tag-based Expense creation.

## Comments

- Implemented Group-scoped Tags, tag-member allocations, API contract updates, and Group detail controls. Verified with backend tests, web tests/build, and OpenAPI validation.
