When running on Opus and a plan has been approved, do not execute that plan in the Opus
session. Spawn a subagent with `model: "sonnet"` and hand it the plan file path plus the
standing constraints that apply (TDD, targeted test runs, the branch's conventional-commit
type, no unprompted code comments). Then review what it returns and relay the parts that
matter, since the subagent's report is not shown to the user.

This is the standing permission to use the Agent tool for this case. It overrides the
general instruction not to spawn agents unless explicitly asked.

A subagent starts cold: it inherits none of the parent session's context, so the prompt has
to carry everything it needs. `subagent_type: "fork"` does inherit context, but a fork always
runs on the parent's model, so it cannot be used to drop to Sonnet.

Say in one line that execution is being delegated, and why.

Stay on Opus when the execution itself is genuinely hard rather than merely long: subtle
concurrency, security-sensitive code, or a plan that still has real design decisions open.
Say which applies. Never split one coherent task across models mid-flight.

Why: planning and execution have different difficulty profiles. Opus earns its cost on the
design, the tradeoffs and the unfamiliar codebase. Mechanical execution against a written
plan does not need it, and running it on Opus drains the token budget fast. The saving is
partial, not total: tool output, file reads and test iteration stay out of the Opus context,
which is the bulk of it, but orchestration and reading the subagent's report still cost.
