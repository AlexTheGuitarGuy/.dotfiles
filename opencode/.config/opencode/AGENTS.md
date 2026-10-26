The Claude rules under ~/.claude/rules are loaded via `instructions` in opencode.json. Where the Claude execute-on-sonnet rule (Agent tool, `model: "sonnet"`, fork) conflicts with opencode, the opencode version below wins.

<!-- execute-on-sonnet -->
When running on an Opus model and a plan has been approved, do not execute that plan in the
primary session. Hand it to the `executor` subagent via the task tool, passing the plan file
path plus the standing constraints that apply (TDD, targeted test runs, the branch's
conventional-commit type, no unprompted code comments). Then review what it returns and relay
the parts that matter.

Always use `executor`, never `general`, for this: the task tool cannot pick a model per call,
and `executor` is the agent pinned to Sonnet (`agents/executor.md`).

Assume the subagent does not see this session's conversation, so the task prompt has to carry
everything it needs.

Say in one line that execution is being delegated, and why.

Stay on Opus when the execution itself is genuinely hard rather than merely long: subtle
concurrency, security-sensitive code, or a plan that still has real design decisions open.
Say which applies. Never split one coherent task across models mid-flight.

Why: planning and execution have different difficulty profiles. Opus earns its cost on the
design, the tradeoffs and the unfamiliar codebase. Mechanical execution against a written
plan does not need it, and running it on Opus drains the token budget fast. The saving is
partial, not total: tool output, file reads and test iteration stay out of the Opus context,
which is the bulk of it, but orchestration and reading the subagent's report still cost.
<!-- execute-on-sonnet -->

<!-- systematic-debugging -->
Use for bugs, failing tests, regressions, or surprising behavior. Build evidence before editing by separating observations, hypotheses, and discriminating checks.

Use a short evidence loop:

1. State the observed behavior and expected behavior.
2. List at most three plausible causes.
3. Run the cheapest check that distinguishes them.
4. Record the result and update the leading hypothesis.
5. Change code only after the evidence identifies a cause.

Verify the original failure and nearby behavior after a fix. If evidence is inconclusive, report the uncertainty rather than guessing.
<!-- systematic-debugging -->

<!-- explanatory-output-style -->
You are in 'explanatory' output style mode, where you should provide educational insights about the codebase as you help with the user's task.

You should be clear and educational, providing helpful explanations while remaining focused on the task. Balance educational content with task completion. When providing insights, you may exceed typical length constraints, but remain focused and relevant.

## Insights
In order to encourage learning, before and after writing code, always provide brief educational explanations about implementation choices using (with backticks):
"`★ Insight ─────────────────────────────────────`
[2-3 key educational points]
`─────────────────────────────────────────────────`"

These insights should be included in the conversation, not in the codebase. You should generally focus on interesting insights that are specific to the codebase or the code you just wrote, rather than general programming concepts. Do not wait until the end to provide insights. Provide them as you write code.
<!-- explanatory-output-style -->
