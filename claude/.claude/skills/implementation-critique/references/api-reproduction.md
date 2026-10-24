# Clean curl reproductions

Prepare requests the human can run against their existing local app. Keep the flow visible: initial state → action → response/state to inspect.

For bug-fix reviews, include the minimal verification requests even when there are no findings. Show the tested application request, necessary state setup, observed status and the response fields that establish the expected business result. A successful direct provider comparison supplements application verification; identify which boundary each request exercises. Convert requests executed through SDKs or temporary helpers into equivalent standalone cURLs for the handoff.

## Ground the scenario

Derive methods, full routes, authentication, request fields, and necessary data from route registration, schemas, handlers, or existing tests. State the minimum precondition in one or two sentences. Use the configured local address when known; otherwise mark the port and identifiers for replacement. Never invent an endpoint, credential, or existing record.

Keep values responsible for the bug concrete, such as a 70/30 ownership split with 50/50 savings. Prefer the tested URL and exact test user and record IDs. Keep secrets such as `<token>` as placeholders; use other placeholders only for unavailable values and explain how to obtain them. Test data follows the test-data rule in `~/.claude/rules/`. Do not invent a fixture server to make the example work.

## Request format

Use one standalone `curl` per fenced block, in execution order. Prefer explicit `--request`, `--url`, `--header`, and inline `--data-raw` JSON. This keeps the request readable and easy to copy into a request client. For example (adapt the route and body to the actual contract):

```sh
curl --request POST \
  --url 'http://localhost:<port>/v1/widgets' \
  --header 'Content-Type: application/json' \
  --data-raw '{
    "name": "reproduction-example"
  }'
```

- Put the step name outside the code block: “1. Save the allocation”, “2. Preview the same owner”.
- Include only required headers and fields. Add `Authorization: Bearer <token>` only when that scheme is established by the source.
- Use literal URLs and inline payloads. Avoid shell variable setup, command substitution, heredocs, loops, `jq`, file payloads, helper scripts, and output-capture plumbing.
- When step two needs an ID from step one, say which response field to copy and include the observed test ID when available; explain that rerunning creation returns a new ID. Otherwise use a clearly named placeholder.
- Include a read request when the result is persisted state rather than the immediate response. Show the relevant status/error code/field beside the requests; a short response excerpt is enough. Add `--include` only when terminal status/header visibility is needed.
- Keep real secrets and captured private payloads out of requests. Do not disable TLS verification.

## Expected and actual behavior

State the intended result and the discrepant result precisely. Label the latter **observed locally** or **observed in <test environment>** when executed, or **from code; not run** when inferred. Include the distinguishing HTTP status, error code or value so an unrelated failure cannot be mistaken for this issue.

If a precondition requires historical database state that today's API cannot create, say so. Point to the existing seed or provide a small, schema-grounded local SQL snippet when needed; do not claim curls alone can establish that state. If the necessary local data is unavailable, state that limitation and finish the finding.

Run the minimal sequence for bug-fix reviews; for other HTTP findings, prefer live execution, using the repo's existing request tooling where it exists. Test data and identities follow the test-data rule in `~/.claude/rules/`. Describe data changes and report observed responses separately from expected ones. If execution remains blocked, identify the actual blocker and label the requests unexecuted. Keep any temporary helpers proportional to the reproduction. Cleanup instructions should target only records created by the reproduction, where supported.
