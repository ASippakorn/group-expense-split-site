# 02 - Create Expense Ledger with Equal Splits

**What to build:** Participants can create an Expense with amount, Payer, description, date, currency, selected Participants, and an equal Split. The Group detail view shows the Expense list and each Expense's Split details.

**Blocked by:** 01 - Add Participant Management.

**Status:** resolved

- [x] A Participant can create an Expense for a Group they participate in.
- [x] The Payer must be one of the Expense Participants.
- [x] Amounts are stored as integer minor units.
- [x] THB is the default currency, while the original Expense currency is stored.
- [x] Equal Split rounding is deterministic.
- [x] The Group detail view lists Expenses with Payer, amount, description, date, currency, and Split details.
- [x] OpenAPI documents the Expense creation and listing routes.
- [x] Backend tests cover equal Split creation, Payer validation, Participant membership validation, and rounding.
- [x] Frontend tests cover creating an equal Split Expense and rendering it in the Group detail view.
