---
name: implementation-critique
description: Review a concrete diff, commit, branch range, or path set against personal standards for locality, misleading intent, type modeling, and shotgun surgery. Return short findings with plain scenarios, expected versus actual behavior, and clean curl requests for local reproduction.
---

# Implementation Critique

Review only an explicit change set. Ask for a range or path set when the scope cannot be inferred safely. Findings advise the reviewer; they do not create a workflow gate.

## Keep the review lightweight

The default deliverable is a short review with a plain sequence of user actions and copyable requests. For bug-fix branches, verifying the reported failure is part of the review: execute the minimal reproduction against the available local/test environment before concluding, with concrete tested IDs and observed results. Do not wait for the user to ask whether the fix actually works.

Test data and identities: follow the test-data rule in `~/.claude/rules/` when one applies.

Run only targeted checks per `targeted-test-runs.md`. For deployed environments use the `review-live` skill. Reproduce API findings with the repo's existing request tooling where it exists; curl blocks are handoff text. Keep production source unchanged and leave publishing comments to the human unless requested.

When the PR belongs to someone else, `review-others-pr.md` wins: no fixes, no new scripts, no screenshots.

Keep verification proportional to the finding. Prefer existing tooling and focused requests or tests; add a small temporary helper when useful. Build indexes or compare another revision when the evidence calls for it. Avoid turning a reproduction into a separate test-platform project or generating evidence bundles the user did not request.

## Gather bounded evidence

1. Resolve the changed files and hunks with Git. Record the review base, head, and relevant working-tree changes; a branch name alone does not identify the reviewed code.
2. Fetch the related ticket before judging requirement compliance. Discover its key or URL from the user's request, branch name, commit messages, and PR title/body; do not wait for the user to paste it. Read its description, acceptance criteria, relevant fields and all comment pages. Use the available authenticated tracker integration, or the ticket exports a rule describes.
   Follow relevant parents/epics, subtasks, linked or referenced tickets, linked design documents and PR discussion when they could clarify intent, dependencies or affected behavior. Read the comments on each ticket consulted. If the code or discussion suggests another ticket is involved without an explicit link, perform a focused read-only search within the relevant project/application and verify the relationship. No extra approval is needed for this contextual discovery. Stop expanding when the requirements and dependencies needed to assess the change are understood; avoid unrelated tracker exploration.
   Cite ticket keys or URLs for requirements used in findings. Treat newer comments as possible clarifications, and surface contradictions instead of silently choosing a requirement. If no related ticket can be identified or access/comment pagination fails, report that limitation and continue the source review without claiming ticket compliance; ask for missing context only where it changes a conclusion. Retrieval does not authorize tracker writes.
3. Read each changed hunk in the context of its owning symbol.
4. Trace immediate callers, consumers, or schema boundaries only where a finding depends on them. Prefer bounded searches and reuse existing indexes when useful; avoid whole-repository orientation for a localized finding. Compare behavior with the requirements and design context gathered above.
5. When the change set affects an HTTP API, trace the route, method, authentication shape, request schema, and required state for each potential API-observable finding.
6. Establish the expected behavior from a named requirement, existing contract, or explicit invariant. Do not use the current implementation as its own specification. Where product intent is missing, report a question.

Use the checks below as the personal rubric and apply concrete additional rules supplied by the user or repository. Tie each finding to the rule it violates; do not turn a preference into a correctness claim without a demonstrated invariant.

Do not equate a large diff or a cross-layer feature with bad design. Separate:

- feature-local changes required for the behavior,
- shared changes required by a demonstrated dependency,
- incidental changes with no demonstrated path from the feature.

## Checks

### Shotgun surgery

Look for one conceptual rule repeated or manually synchronized across unrelated files. A valid finding must cite the affected changed files and name a concrete owner or consolidation point.

Do not report necessary controller/domain/persistence/client wiring, generated artifacts, or call-site updates for an intentionally changed public contract merely because several files changed.

### Locality of behavior

Look for business decisions living in transport, presentation, fixtures, scripts, package commands, or infrastructure instead of beside the state or invariant they govern. Also look for behavior fragmented into helpers whose only effect is to make the owning rule harder to see.

Before reporting, verify that the apparent layer violation is not an adapter-specific translation or boundary concern.

### Misleading intent

Compare concrete claims with behavior:

- names and signatures versus function bodies,
- commit or task scope versus the actual concerns changed,
- comments, tests, and design-document literals versus implementation,
- typed domain errors versus their observable transport behavior,
- apparently passive operations that calculate, default, discard, or throw.

Report only a specific contradiction. Do not infer product intent that is absent from the evidence; ask a question instead.

### Type modeling

Look for observable domain facts weakened at internal boundaries:

- branded or validated values converted back to primitives and reparsed,
- `unknown`, `any`, booleans, optionals, or string codes hiding a finite domain state,
- invalid combinations representable where a discriminated union could prevent them,
- generic errors replacing an existing typed error model,
- schema widening or strictness changes with unexamined consumers,
- duplicated validation caused by unclear boundary ownership.

Prefer the smallest type that makes the demonstrated invariant explicit. Do not demand brands, unions, or abstractions for values without a real domain distinction.

## Verification

For a bug-fix branch, turn the ticket's reported failure and acceptance criteria into the smallest meaningful verification sequence. Start with the reported identity, records, settings and actions when available; use equivalent synthetic test data when necessary and state that substitution. A ticket's successful reproduction is prior evidence, not verification of the reviewed revision. Matching the requested code change or passing tests with mocked providers does not establish that the reported failure is resolved.

Run the exact reviewed commit, and verify the active checkout/build before attributing results to it. A newer deployed version or a different working tree cannot verify this branch. If a provider is involved, distinguish mapper/unit checks, a direct live provider request, and the request through the application. Follow the failing path through the application to its user-visible result; provider success alone leaves application verification incomplete.

Use a controlled before/after comparison where feasible, keeping data, identity and configuration equivalent. Running the base revision is strongest; replaying the failing provider payload with only the changed field restored can isolate the mechanism, but must be labeled as a payload comparison rather than a base-branch application run. Check the returned business result and required preserved values, not just a successful HTTP status. Cover acceptance-criteria variants and nearby state transitions that the change touches, such as fresh versus persisted selection or empty versus populated input, without expanding into an exhaustive suite.

Resolve configuration or fixture failures through the documented setup before interpreting them as defects. Confirm whether focused tests use snapshots/mocks or live providers, and choose the intended mode explicitly when needed. When interrupted and resumed, reuse completed evidence and check running processes and partial mutations before restarting work. Restore temporarily changed shared test state where supported and stop only the verification services you started; report retained fixtures or runtime changes.

For each behavioral finding, state the required initial data, the actions to perform, and the exact response or state to inspect. For HTTP findings, read [the curl reproduction guide](references/api-reproduction.md) and include the complete minimal requests directly in the finding. The user should be able to replace clearly marked local values and run the sequence without a generated script.

Label actual behavior **observed locally** or **observed in <test environment>** only when executed, with the revision when known and any material mocked boundary. Otherwise label it **from code; not run** and give the source-supported prediction. A setup failure or unrelated error does not confirm the finding. Do not claim a provider calculation was verified when only its outgoing request was inspected.

If HTTP cannot expose the issue, use the nearest existing focused test or precise source/state check. Architectural findings need concrete dependencies and change cost, not invented curls or failing tests. For bug-fix reviews, continue until the original failure and relevant acceptance criteria are verified or a concrete blocker remains after available setup and diagnostic checks. A blocked review can report no source findings, but must say the fix is not fully verified and identify the untested boundary. Do not stop at a disclaimer while an authorized verification step remains available. Exhaustive testing or automation of every case is not required.

## Finding threshold

Every finding must include:

- a changed path and symbol or exact hunk,
- the demonstrated behavior or dependency,
- the violated invariant or future change cost,
- one focused remediation direction.

Severity:

- **critical**: proven silent corruption, security issue, or broken public/domain invariant.
- **major**: proven recurring change scattering, misleading public behavior, or invalid-state modeling at a material boundary.
- **minor**: localized maintainability problem with concrete evidence.

Do not emit style-only findings by default. Cap severity at the strongest invariant actually demonstrated. Mark incomplete evidence as an open question rather than upgrading speculation into a finding.

## Output

Include the reviewed range once, then one self-contained Markdown finding per issue using [the finding template](assets/finding.md): scenario, expected behavior, actual behavior, clean requests, and a short source/cause note. The finding itself should be suitable for an author thread; do not append a duplicate comment or lengthy review report.

Return findings in the response by default. If asked to save Markdown files, use the requested or existing review directory and write one file per issue. No generated index, runner, or evidence directory is needed. Keep prose to a few short paragraphs; include only the request fields and response excerpt needed to reproduce and understand the problem.

Mention material uncertainties briefly. If no findings survive, say so and summarize the relevant coverage. For a bug-fix review, report separately whether source review found issues and whether the original failure was verified as resolved, naming the tested commit and environment. Include the minimal copyable application requests and observed status/response excerpt even when the fix passes and there are no findings; a results table alone is not a runnable handoff. State material substitutions, mocked boundaries and remaining gaps without overstating confidence. Do not fill the report with mandatory blast-radius, confidence, non-findings, or remediation-plan sections.
