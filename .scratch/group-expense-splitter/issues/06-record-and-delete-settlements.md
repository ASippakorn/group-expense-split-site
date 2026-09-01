# 06 - Record and Delete Settlements

**What to build:** Participants can record trusted Settlements, delete and recreate them, and see Balances update from Expenses, Splits, and Settlements.

**Blocked by:** 03 - Calculate Group Balances.

**Status:** ready-for-review

- [x] Either side of a repayment can record a Settlement.
- [x] A Settlement records payer, receiver, amount, currency, date, and note.
- [x] Settlement amounts may exceed the currently owed amount.
- [x] Settlements affect Balances immediately.
- [x] A Settlement can be deleted but not edited in place.
- [x] Deleted Settlements no longer affect Balances.
- [x] The Group detail view lists Settlements and supports recording and deleting them.
- [x] OpenAPI documents Settlement routes.
- [x] Backend tests cover recording, overpayment, deletion, Balance changes, and authorization.
- [x] Frontend tests cover Settlement form success, deletion, and Balance update display.

## Comments

- Implemented Settlement persistence, authenticated API routes, derived Balance updates, Group-detail form/list/delete controls, OpenAPI documentation, and focused service/frontend coverage. Awaiting final review.
