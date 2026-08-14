# 04 - Support Manual Amount and Percentage Splits

**What to build:** Participants can create manual amount and percentage Splits, with backend validation that totals exactly match the Expense amount or 100 percent.

**Blocked by:** 03 - Calculate Group Balances.

**Status:** ready-for-agent

- [ ] A Participant can create a manual amount Split.
- [ ] Manual amount Splits are rejected unless the total exactly matches the Expense amount after rounding.
- [ ] A Participant can create a percentage Split.
- [ ] Percentage Splits are rejected unless the total is exactly 100 percent.
- [ ] Balances update correctly for manual amount and percentage Splits.
- [ ] The Group detail view displays the Split type and Participant-level Split values.
- [ ] OpenAPI documents equal, manual amount, and percentage Split inputs.
- [ ] Backend tests cover valid and invalid manual amount and percentage Splits.
- [ ] Frontend tests cover switching Split types and showing validation errors.
