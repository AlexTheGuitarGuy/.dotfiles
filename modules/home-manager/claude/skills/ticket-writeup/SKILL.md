---
name: ticket-writeup
description: Format concise, ticket-ready reports for bugs, issues, behavior changes, and feature implementations. Use when reporting a problem, describing requested behavior, or handing off an implemented change; contrast current and expected behavior and include copyable curls with expected HTTP status and response bodies for API work. Do not use it to decompose a goal or publish tracker tickets.
---

# Ticket writeup

Use this skill alongside the investigation, review, or implementation workflow that establishes the facts. It controls the shape of the report only; it does not authorize code changes or ticket publication.

## Ground the behavior

- Describe concrete externally visible behavior, not a general problem statement.
- Establish current behavior from inspected source, an existing contract, or an executed reproduction. Label it **observed locally** or **observed in <test environment>** only when it was run; otherwise say **from code; not run**. Do not present an assumption as current behavior.
- Establish expected behavior from the user's request, a named requirement, or an existing invariant. State a specific status, response field, value, side effect, or visible state.
- For a feature or completed implementation, keep the pre-change state under **Current behavior** and the intended contract under **Expected behavior**. Add a one-line status saying whether the expected behavior is proposed, implemented, or verified.
- Do not invent routes, payloads, credentials, responses, product rules, or reproduction results. State the exact missing fact when it prevents a concrete example.

## Write the report

Use this order:

```markdown
# <specific problem or outcome>

**Status:** <proposed | implemented | verified>

## Current behavior

<What happens before the change, including the triggering state or input. Say whether this was observed or derived from code.>

## Expected behavior

<What should happen for the same state or input.>

## API example

<Include only for HTTP/API behavior.>
```

Keep each section to the minimum needed to distinguish the behaviors. Prefer concrete example values. Omit background, root-cause speculation, raw logs, exhaustive acceptance criteria, estimates, rollout plans, and implementation detail unless the user asks for them or they are necessary to understand the behavior. If a cause or implementation note is material, add one short paragraph after the behavior contract.

For a non-API issue, replace **API example** with a short reproduction only when the trigger is not already clear from **Current behavior**.

## API examples

For every API flow:

Prefer a live minimal reproduction through the accompanying review/debug workflow before writing the cURLs. Test data and identities follow the test-data rule in `~/.claude/rules/`. If execution is blocked, state the missing prerequisite and mark the requests unexecuted; never invent observed results.

1. Show the smallest standalone `curl` that demonstrates the behavior. Use the evidenced method and full route, necessary headers, and an inline JSON request body.
2. Use the tested URL and exact test user and record IDs when available. Keep credentials such as `{{token}}` as placeholders; use other placeholders only for values that remain unavailable and explain how to obtain them.
3. Immediately show the expected HTTP status and the smallest valid JSON response body that proves the contract. When the response is `204`, state `no response body`.
4. For a bug, also show the current status/body when it materially differs and is known. Label it **observed locally**, **observed in <test environment>**, or **from code; not run**.
5. For a multi-request flow, number requests in execution order and say which returned field is reused. Include the expected response for each request.

Use a copyable request shape:

```bash
curl --request POST \
  --url '{{base_url}}/v1/widgets' \
  --header 'Authorization: Bearer {{token}}' \
  --header 'Content-Type: application/json' \
  --data-raw '{
    "name": "example"
  }'
```

**Expected response:** HTTP 201

```json
{
  "id": "<generated-id>",
  "name": "example"
}
```

Replace the illustrative route and fields with the real contract. Do not add token-acquisition commands, shell variables, command substitution, `jq`, loops, retries, setup scripts, or unrelated response fields. Never expose secrets or captured private data.


## Posting

Text that goes to Jira or a PR follows `external-writes.md`: it reads as the user's own voice and uses no dashes. Posting needs the user's approval.
