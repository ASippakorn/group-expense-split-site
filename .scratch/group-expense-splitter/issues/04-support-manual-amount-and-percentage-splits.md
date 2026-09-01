# 04 - Support Manual Amount and Percentage Splits

**What to build:** Participants can create manual amount and percentage Splits, with backend validation that totals exactly match the Expense amount or 100 percent.

**Blocked by:** 03 - Calculate Group Balances.

**Status:** resolved

- [x] A Participant can create a manual amount Split.
- [x] Manual amount Splits are rejected unless the total exactly matches the Expense amount after rounding.
- [x] A Participant can create a percentage Split.
- [x] Percentage Splits are rejected unless the total is exactly 100 percent.
- [x] Balances update correctly for manual amount and percentage Splits.
- [x] The Group detail view displays the Split type and Participant-level Split values.
- [x] OpenAPI documents equal, manual amount, and percentage Split inputs.
- [x] Backend tests cover valid and invalid manual amount and percentage Splits.
- [x] Frontend tests cover switching Split types and showing validation errors.

## Comments

- Implemented manual amount and percentage Split creation, persistence, validation, API documentation, and Group-detail controls. Percentage allocations use deterministic largest-remainder rounding in minor units.
- Verified with `go test ./...`, web tests and typechecking, and OpenAPI validation. Standards and spec review completed; the one display gap found during review was corrected.
