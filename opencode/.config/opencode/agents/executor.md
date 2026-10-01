---
description: Executes an already-approved plan file step by step. Use only to hand off an approved plan from an Opus session.
mode: subagent
model: github-copilot/claude-sonnet-5-5
---
Execute the approved plan at the path given in the task prompt, following the constraints
passed with it. Do not revisit design decisions the plan already made. If the plan is wrong
or blocked, stop and report why instead of improvising.
