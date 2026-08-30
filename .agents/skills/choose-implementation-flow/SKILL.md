---
name: choose-implementation-flow
description: "Assess a local issue before implementation and recommend /tdd or /implement, model tier, reasoning effort, and review depth."
disable-model-invocation: true
---

# Choose Implementation Flow

Use this before implementing a local issue or ticket when the user wants to spend model tokens deliberately.

## Process

1. **Read the issue.** If the user gives an issue number, find the matching file under `.scratch/*/issues/`. Read the full issue, its `Status:`, `Blocked by:`, and checklist.
2. **Check blockers.** Read any local blocker issues named by `Blocked by:`. If a blocker is not resolved, recommend waiting or implementing the blocker first.
3. **Read the feature spec and domain context.** Read the nearest `.scratch/<feature>/spec.md` when present, plus `CONTEXT.md` when present. Use them to judge the ticket's breadth and risk.
4. **Lightly inspect the codebase.** Look only far enough to identify touched surfaces: backend, frontend, database, OpenAPI, tests, CI, infrastructure, or security-sensitive areas.
5. **Classify the work.**
   - **Tiny:** one small behavior or UI change, one surface, low risk.
   - **Focused:** one main behavior with clear tests, one or two surfaces.
   - **Ticket-sized:** several connected checklist items, cross-surface work, API/UI/test coordination, or persistence.
   - **Risky:** auth, permissions, money, data integrity, migrations, concurrency, security, large refactors, or unclear architecture.
6. **Recommend the flow.**
   - Use `/tdd` for Tiny or Focused work where one public seam can drive the change.
   - Use `/implement` for Ticket-sized work or any issue with multiple connected checklist items.
   - Use `/diagnosing-bugs` first when the issue is a hard bug without a reliable failing command.
   - Use `/codebase-design` or `/prototype` first when the main uncertainty is where the interface should live or whether a state model/UI will feel right.
7. **Recommend the model and effort.**
   - **Tiny:** `gpt-5.6-luna`, `low` or `medium`.
   - **Focused:** `gpt-5.6-terra`, `low` or `medium`.
   - **Ticket-sized:** `gpt-5.6-terra`, `medium`.
   - **Risky:** `gpt-5.6-sol`, `high`; use `xhigh` only for hard architecture, data correctness, or repeated failed attempts.
8. **Recommend review depth.**
   - **Light review:** main-agent self-review against the issue checklist, spec, tests, and nearby conventions.
   - **Full `/code-review`:** use when the user asks, the diff is large, or the change touches auth, permissions, money, data integrity, migrations, security, shared architecture, or production/deployment behavior.

## Output

Return a compact recommendation:

- **Flow:** `/tdd`, `/implement`, or another prerequisite skill.
- **Model:** exact model and reasoning effort.
- **Review:** light/self-review or full `/code-review`.
- **Why:** two or three concrete reasons from the issue and codebase.
- **Suggested prompt:** a ready-to-run sentence the user can paste into the next implementation session.

If the issue is blocked, put the blocker first and do not recommend spending implementation tokens yet.
