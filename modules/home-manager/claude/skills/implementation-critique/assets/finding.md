# <severity>: <concrete problem>

**Scenario:** <minimum initial state and ordered user actions; include the values that cause the issue>.

**Expected:** <specific status, value, or state>.

**Actual (<observed locally at revision / observed in test environment / from code; not run>):** <specific discrepant result and relevant error code or value>.

## Reproduce locally

<Only the necessary prerequisites and placeholder replacements.>

1. <First action.>

```sh
<standalone curl with explicit URL, necessary headers, and inline body>
```

2. <Next action or observation; say which returned ID to reuse if needed.>

```sh
<standalone curl>
```

**Cause:** [<path and symbol>](<source-location>): <violated review rule, why the behavior occurs, and focused correction in one or two sentences>.

<!-- Replace placeholders and use only as many requests as needed. For non-HTTP issues, replace curls with an existing focused check or source evidence. The note itself is the author comment; omit duplicate summaries, raw logs, and this template instruction. -->
