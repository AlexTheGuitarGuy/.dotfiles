---
name: systematic-debugging
description: Use for bugs, failing tests, regressions, or surprising behavior. Build evidence before editing by separating observations, hypotheses, and discriminating checks.
---

# Systematic Debugging

Use a short evidence loop:

1. State the observed behavior and expected behavior.
2. List at most three plausible causes.
3. Run the cheapest check that distinguishes them.
4. Record the result and update the leading hypothesis.
5. Change code only after the evidence identifies a cause.

Verify the original failure and nearby behavior after a fix. If evidence is inconclusive, report the uncertainty rather than guessing.
