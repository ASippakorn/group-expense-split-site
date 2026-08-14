# 02 - Create Expense Ledger with Equal Splits

**What to build:** Participants can create an Expense with amount, Payer, description, date, currency, selected Participants, and an equal Split. The Group detail view shows the Expense list and each Expense's Split details.

**Blocked by:** 01 - Add Participant Management.

**Status:** ready-for-agent

- [ ] A Participant can create an Expense for a Group they participate in.
- [ ] The Payer must be one of the Expense Participants.
- [ ] Amounts are stored as integer minor units.
- [ ] THB is the default currency, while the original Expense currency is stored.
- [ ] Equal Split rounding is deterministic.
- [ ] The Group detail view lists Expenses with Payer, amount, description, date, currency, and Split details.
- [ ] OpenAPI documents the Expense creation and listing routes.
- [ ] Backend tests cover equal Split creation, Payer validation, Participant membership validation, and rounding.
- [ ] Frontend tests cover creating an equal Split Expense and rendering it in the Group detail view.
