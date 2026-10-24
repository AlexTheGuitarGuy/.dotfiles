---
name: investigate-code
description: Explain existing code or a diff through a source-grounded walkthrough. For API changes, open with request and expected response examples and the changed behavior flow, then walk through the code bottom-up. Markdown in chat by default; the evidence terminal viewer when asked.
---

# Investigate code

Reverse-engineer the code that exists; do not design a future implementation or modify source. This is a guided explanation for the human, not a findings-only review or a planning workflow.

## Scope and discovery

- Use the explicit files, directory, behavior, commit, or diff range. With no target, use the code or changes under discussion; ask one scope question if that is ambiguous.
- For changes, inspect the actual diff and enough surrounding code to explain it. Keep staged, unstaged, and committed changes distinct. Do not checkout, reset, or discard work.
- Locate the relevant symbols with GitNexus (see gitnexus-exploring: `query`, `context`, `impact`), then read exact source ranges. Report unavailable tooling and use targeted reads rather than inventing relationships.
- For a large scope, map the meaningful layers first and divide them into manageable connected steps. Explicitly report omitted areas; do not silently reduce a large change to one representative function.

## Bottom-up walkthrough

For an API PR or diff affecting HTTP behavior, first prepare an API overview using [the API overview guide](references/api-overview.md). Explain the observable change and intended flow, including concrete request bodies and expected response bodies/statuses, before walking through the implementation. Detect this from the inspected routes, schemas, handlers, and their immediate consumers; do not require a PR lookup just to identify API code.

1. Start with leaf values, data structures, types, invariants, and primitives.
2. Explain the functions, adapters, and domain rules that consume them.
3. Follow callers and data flow upward to orchestration, entry points, and observable behavior.

Each numbered step needs focused, verbatim source captures with paths and line ranges. Briefly explain what each capture does, why it matters, and how its inputs/outputs connect to adjacent steps. Keep code central and prose concise. Explain errors, side effects, and boundary conditions where they occur. For diffs, distinguish existing behavior from changed behavior and explain the consequence.

Use real dependencies, not alphabetical file order. Group cycles into a single step or explicitly explain the cycle. Separate observed facts from inferred intent. Attach concrete concerns and missing verification to the relevant step; do not manufacture findings or claim that reading code proves correctness. End with the complete behavior and remaining uncertainties.

## Presentation

Default to Markdown in the conversation (see the last paragraph). Use the viewer only when the user asks for it.

Viewer mode: write the walkthrough to `.evidence/investigations/walkthrough-<short-slug>/` at the repository root, and add `.evidence/` to `.git/info/exclude` if it is not there. Write only the walkthrough artifacts there. Create `manifest.json` in this shape (fields are top-level, not nested under `payload`):

```json
{
  "mode": "learn",
  "version": 1,
  "goal": "The requested code or change",
  "scope": "Exact files/range and any exclusions",
  "overview": {"summary": "Short end-to-end explanation", "sections": [], "api": []},
  "breakpoints": [
    {
      "id": "leaf-contract",
      "order": 1,
      "title": "Leaf contract",
      "paths": ["src/example.ts"],
      "reason": "Why this step matters",
      "sketch": "What it does and how it feeds the next step",
      "dependsOn": [],
      "downstream": [],
      "sections": [],
      "captures": [
        {
          "path": "src/example.ts",
          "startLine": 1,
          "endLine": 3,
          "label": "Contract",
          "snippet": "Verbatim source at these lines, not illustrative code",
          "reason": "Brief explanation tied to this capture"
        }
      ]
    }
  ],
  "skillReports": []
}
```

- Replace the example with actual inspected code. Every step needs at least one capture with exact positive line numbers. Do not put ellipses or invented code in captures.
- IDs must be unique and use only letters, digits, underscores, or hyphens. `order` is bottom-up; `dependsOn` and `downstream` reference other step IDs.
- Optional semantic sections use `{kind, title, summary, points}` and support the captures rather than replacing them. Do not generate visuals or require diagram renderers.
- `paths`, `dependsOn`, `downstream`, and section `points` are arrays of strings. Source line numbers, step `order`, and API `responseStatus` are integers; API bodies remain actual JSON values.
- Leave `skillReports` empty or omit it unless reporting a specific skill's coverage. Each entry must be a report object, never a string. Put methodology notes, unavailable tooling, and scope limitations in `overview.sections`, using `summary` and string-array `points`. A report uses this shape, with `coverage`, `findings`, and `nonFindings` optional:

  ```json
  {
    "id": "investigate-code",
    "title": "Code walkthrough",
    "status": "checked",
    "coverage": ["Inspected the request handler and its direct dependencies"],
    "findings": [],
    "nonFindings": ["HTTP examples are expected contracts; requests were not executed"]
  }
  ```

  Every finding is an object with string `title`, `evidence`, and `impact`, plus optional string `threadId` and `suggestedComment`. Do not substitute prose strings for report or finding objects, or invent completed checks to populate a report.
- Every new walkthrough opens on a visible **Main overview** page. For API changes, its short sections explain **Before**, **Now**, and the ordered **Flow**; fill `overview.api` with the source-grounded examples described in the guide. The sidebar has two groups: **Overview & examples** (main overview, then one page per API exchange) and **Code bottom-up** (numbered source steps). `j`/`k` navigate all pages, `o` returns to the overview, and `c` jumps to code. Reopening preserves the saved reading position; the overview remains selectable in the sidebar.
- After saving the complete manifest, run the bundled validator, resolving `<skill-directory>` to the directory containing this `SKILL.md`:

  ```sh
  node "<skill-directory>/scripts/validate-manifest.mjs" "<artifact-directory>"
  ```

  It checks field types, report objects, source ranges, unique step IDs, and dependency references. JSON parsing alone does not check the viewer's data contract. Repair any reported errors and rerun until it exits successfully; still verify source fidelity and API claims against inspected code. The validator only reads the artifact and never launches the viewer or changes reading state.
- After validation succeeds, check `command -v evidence`. If it is missing, say so and give the Markdown walkthrough instead. Otherwise open the viewer in its own terminal window as the last action:

  ```sh
  ghostty --working-directory="$PWD" -e evidence learn "<artifact-directory>" >/dev/null 2>&1 &
  ```
- The viewer is a reader: navigation and read/unread progress are local. Closing it never requests changes, asks the agent a question, or approves implementation. Ignore legacy question/request markers; do not consume them or turn them into prompts. Preserve existing human notes and reading state.
- After opening the viewer, report the artifact directory and the reopen command (`evidence learn <artifact-directory>`) in one line and stop. Do not revise the artifact, reopen the viewer, or change source unless the human sends a new request.

In Markdown mode, present the same overview and numbered walkthrough in ordinary Markdown. For very large scopes, present the step map and the first connected chunk, then let the human choose the next chunk or ask questions. Do not create artifacts unless requested.
