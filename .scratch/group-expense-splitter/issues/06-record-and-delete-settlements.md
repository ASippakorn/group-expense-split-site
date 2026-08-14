# 06 - Record and Delete Settlements

**What to build:** Participants can record trusted Settlements, delete and recreate them, and see Balances update from Expenses, Splits, and Settlements.

**Blocked by:** 03 - Calculate Group Balances.

**Status:** ready-for-agent

- [ ] Either side of a repayment can record a Settlement.
- [ ] A Settlement records payer, receiver, amount, currency, date, and note.
- [ ] Settlement amounts may exceed the currently owed amount.
- [ ] Settlements affect Balances immediately.
- [ ] A Settlement can be deleted but not edited in place.
- [ ] Deleted Settlements no longer affect Balances.
- [ ] The Group detail view lists Settlements and supports recording and deleting them.
- [ ] OpenAPI documents Settlement routes.
- [ ] Backend tests cover recording, overpayment, deletion, Balance changes, and authorization.
- [ ] Frontend tests cover Settlement form success, deletion, and Balance update display.
