# Group Expense Splitter

The Group Expense Splitter tracks shared expenses inside groups and uses those expense records plus settlement records to show who owes whom.

## Language

**User**:
A registered person who can sign in and participate in groups.
_Avoid_: Account, member

**Group**:
A named expense-sharing space, such as a trip, meal, or recurring gathering.
_Avoid_: Trip, event, room

**Default Currency**:
The currency a Group uses when creating new Expenses unless a different currency is selected.
_Avoid_: Base currency, home currency

**Participant**:
A User included in a Group for expense sharing.
_Avoid_: Attendee, member

**Inactive Participant**:
A Participant who no longer takes part in a Group but remains attached to historical Expenses, Splits, Settlements, and Balances.
_Avoid_: Removed user, deleted participant

**Owner**:
The Participant who created a Group and manages its Participants.
_Avoid_: Admin, manager

**Regular Participant**:
A Participant who can add Expenses and Settlements in a Group but cannot manage other Participants.
_Avoid_: Editor, viewer

**Expense**:
A cost paid by one Participant and allocated across one or more Participants.
_Avoid_: Bill, transaction

**Split**:
The allocation of an Expense across Participants, either equally, by manual amount, by percentage, or by tag.
_Avoid_: Share, allocation rule

**Tag**:
A Group-level label used to include only selected Participants in a tagged Split.
_Avoid_: Category, label

**Settlement**:
A payment record showing that one Participant repaid another Participant.
_Avoid_: Payment, transfer

**Summary Export**:
A PDF report for a Group over a date range, including Expenses, Settlements, Balances, and Suggested Transfers.
_Avoid_: Report, statement

**Balance**:
The calculated amount a Participant owes or is owed after Expenses and Settlements are applied.
_Avoid_: Debt, net

**Suggested Transfer**:
A calculated recommendation for one Participant to repay another Participant based on current Balances.
_Avoid_: Smart settlement, payment instruction

**Receipt Attachment**:
A stored file linked to an Expense as supporting evidence.
_Avoid_: Upload, proof

**Demo Data**:
Seeded Users, Groups, Expenses, and Settlements used for local development, testing, and review.
_Avoid_: Sample data, fixtures

**Exchange Rate**:
A configured conversion rate used to compare or report Expenses recorded in different currencies.
_Avoid_: Live rate, FX feed

**PromptPay Receiving ID**:
A User's receiving identifier used to generate a PromptPay QR code for Suggested Transfers.
_Avoid_: PromptPay account, payment QR
