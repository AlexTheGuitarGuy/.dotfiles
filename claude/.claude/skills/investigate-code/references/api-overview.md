# API behavior before the code

Use this for a PR, diff, or selected code that changes an HTTP API's observable behavior. Give the reader the flow they are trying to understand before the bottom-up source walkthrough. Keep it to the affected API behavior; preparing an overview does not require executing requests, creating fixtures, or building a test harness.

## Explain the change

Use `overview.summary` for a short purpose statement naming the affected endpoints. Use `overview.sections` for three concise sections on the main overview page:

- **Before:** how the endpoint behaved at the diff's base, including what was accepted, stored, or returned.
- **Now:** what the change means for the caller; mention unchanged response contracts when only the internal behavior changed.
- **Flow:** the relevant starting state and ordered steps, such as “save allocation → persist both owners → preview with the saved shares → return the allocation.” Use `points` for the steps and keep the same example identifiers across requests.

Explain changed validation, side effects, and error outcomes only where they matter to this flow. Inspect the base revision to establish “Before”; if it is unavailable, state that rather than guessing. Keep this page about behavior, not file-by-file implementation. Put detailed bodies on the separate example pages.

Ground methods, full paths, required inputs, status codes, and response fields in inspected routes, schemas, handlers, serializers, or existing contract tests. Point to those sources in the overview sections or the linked code steps. Separate the expected contract from any discrepancy in the implementation. If intent or a response field cannot be established, say so instead of inventing it.

## Show concrete exchanges

Add one entry to `overview.api` for each request needed to follow the changed behavior. Include a representative success exchange and an error exchange only when that error is part of the change. For a sequence, keep the same example identifiers and make the order clear in titles/flow sections.

```json
{
  "title": "1. Save the allocation",
  "method": "PUT",
  "path": "/allocations/example-id",
  "summary": "Illustrative contract example. Explain required state, authentication, and what this step changes.",
  "requestBody": {"share": 0.7},
  "responseStatus": 200,
  "responseBody": {"id": "example-id", "share": 0.7}
}
```

Replace this illustrative route and payload with the inspected contract. Bodies are actual JSON values, not JSON encoded inside strings; objects, arrays, scalars, and explicit `null` are supported. Omit `requestBody` or `responseBody` only when that exchange has no body, such as a GET request or 204 response. Omit `responseStatus` if it cannot be established and explain the uncertainty in `summary`. Query parameters belong in `path`; mention any required headers in `summary` without including credentials.

Use a complete minimal valid request and the response shape needed to understand the flow. Label excerpts if a large response is shortened. Examples are illustrative expected behavior, not observed HTTP results; do not imply they were run. Use obviously synthetic identifiers and no real tokens or private captured payloads.

In Markdown mode, use the same content as headings: **Before**, **Now**, **Flow** (numbered steps), then one subsection per exchange with method and path, a `json` request block, the status, and a `json` response block.

In viewer mode, the TUI opens on the main behavior overview. Its **Overview & examples** group exposes each exchange as a separate page with formatted request and expected response bodies; **Code bottom-up** contains the numbered source captures explaining the implementation. Non-API walkthroughs omit `overview.api` but still provide a short main overview of the behavior before the code steps.
