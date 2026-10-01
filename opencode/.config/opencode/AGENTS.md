<!-- context7 -->
Use Context7 MCP to fetch current documentation whenever the user asks about a library, framework, SDK, API, CLI tool, or cloud service -- even well-known ones like React, Next.js, Prisma, Express, Tailwind, Django, or Spring Boot. This includes API syntax, configuration, version migration, library-specific debugging, setup instructions, and CLI tool usage. Use even when you think you know the answer -- your training data may not reflect recent changes. Prefer this over web search for library docs.

Do not use for: refactoring, writing scripts from scratch, debugging business logic, code review, or general programming concepts.

## Steps

1. Always start with `resolve-library-id` using the library name and the user's question, unless the user provides an exact library ID in `/org/project` format
2. Pick the best match (ID format: `/org/project`) by: exact name match, description relevance, code snippet count, source reputation (High/Medium preferred), and benchmark score (higher is better). If results don't look right, try alternate names or queries (e.g., "next.js" not "nextjs", or rephrase the question). Use version-specific IDs when the user mentions a version
3. `query-docs` with the selected library ID and the user's full question (not single words)
4. Answer using the fetched docs
<!-- context7 -->

<!-- no-em-dash -->
Never use the em dash (—) or en dash (–) in written text: chat responses, code comments, commit messages, documentation, or any other prose, in any project.

Use a comma, period, semicolon, parentheses, or a connecting word (and, but, since, because) instead.
<!-- no-em-dash -->

<!-- no-em-dash-scope -->
The no-em-dash rule (see above) governs only text I author myself.

Never retroactively "clean up" an em dash or en dash in existing file content that someone else wrote, even while substantially rewriting other parts of the same file for an unrelated task. When editing a file for a specific task, only avoid the dashes in the new text being added; do not scan the rest of the file for dashes to fix as a side effect. If genuine em-dash cleanup of existing prose is wanted, that must be an explicit, separate request from the user.
<!-- no-em-dash-scope -->

<!-- external-writes -->
Any write or modification to a system visible to other people (Jira tickets, Confluence pages, GitHub/PR comments, Slack messages, emails, or similar) requires explicit approval from the user before posting, every time, no standing permission.

Before posting, the content must:
1. Read as authored by the user, not by an AI. No AI-tell phrasing, no signature or disclosure that it was AI-written.
2. Follow the no-em-dash rule above, since this is user-facing/external-facing text.

Read-only lookups (viewing tickets, pages, threads) need no approval. Any create/update/comment/transition action always needs a review-and-approve step first, drafted in plain prose without dashes or AI-sounding phrasing.

Why: this kind of content is visible to teammates and represents the user's own voice/authorship, unlike local file edits or memory notes.
<!-- external-writes -->

<!-- research-permissions -->
Before doing open-ended research (e.g. "what is X"), across any project:

1. Use any CLI needed (az, kubectl, helm, git, etc.) without asking first, for read-only exploration.
2. Look anywhere in the relevant workspace/repos. If a repo isn't on the latest commit, local changes may be stashed and main checked out and pulled, but ask permission first before that checkout/pull/stash step.
3. Use any MCP needed for the research.
4. Never touch prod or staging/UAT environments without asking first. Any non-readonly operation (write/delete/exec/restart) on a pod, database, or similar resource, in any environment, also needs to be asked about first. Exception: DFS test data created and
   changed through the ai-test-data scripts on d0x/entwicklung/integration (see the ai-test-data rule below).

This is standing permission for read-only exploration (grep, CLI calls, MCP queries) so research can be self-served rather than asking each time, with carve-outs for repo state mutation and for anything touching real infrastructure.
<!-- research-permissions -->

<!-- safe-delete -->
Before running any destructive operation (rm, mv as a replace step, overwrite) on a
path, resolve it first with `readlink -f <path>` or `realpath <path>` and compare
against the literal path given. If they differ, an ancestor directory in the path is
a symlink, not just the leaf, so checking whether the target file itself is a symlink
is not enough. Stop and confirm before proceeding if the resolved path points
somewhere unexpected (e.g. into a git repo, a stow-managed tree, or any location
other than where the literal path suggests).

Never delete-then-recreate when replacing a file or directory. Always copy the real
content to a separate temp location first, verify the copy succeeded, and only then
remove the original. If a later step in the replace sequence fails, the copy still
exists and nothing is lost.

Why: a real incident where `~/.config/zsh` was a whole-directory symlink into a git
repo. The file being deleted was a plain regular file, not a symlink, so a check on
the leaf alone would have missed it. `rm` ran before a copy existed anywhere else,
permanently destroying the only copy of the file when the following `mv` failed.
<!-- safe-delete -->

<!-- verify-claims -->
Use before claiming a change, review finding, plan, diagnosis, or test result is correct. Convert the claim into concrete evidence and record uncertainty.

For each material claim, capture:

- **Claim:** what is believed to be true.
- **Evidence:** path, command output, test result, or direct observation.
- **Scope:** what the evidence does and does not cover.
- **Confidence:** high, medium, or low, with the reason.

Prefer fresh, local evidence. A passing command is not proof beyond what it exercised. Do not turn verification into a stage or a prerequisite; emit it as evidence for later work.
<!-- verify-claims -->

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

<!-- commit-type-matches-branch -->
Every commit on a branch must use the branch/PR's conventional-commit type, not the literal nature of that individual change. If the branch/PR is `feat(...)`, every commit on it is `feat(...)` too, even ones that would look like a `fix` or `docs` commit in isolation (bug fix, doc update, test-only change). Never mix `fix(...)`, `docs(...)`, etc. into a `feat(...)` branch, or vice versa.

Before writing any commit message, check the branch name, PR title, or the branch's first commit for its conventional-commit type, and prefix every subsequent commit on that branch with that same type.
<!-- commit-type-matches-branch -->

<!-- morning-brief -->
Use for daily triage, "what needs attention", or a compact current-work GitHub summary — active PRs authored by the user, PRs awaiting their review, and unread notifications. Trigger on explicit ask, never automatically at session start.

Steps:
1. `gh pr list --author=@me --state=open` for PRs the user has open, scoped to `dvag/*` orgs by default (matches the repo mapping already in `gh-dash/.config/gh-dash/config.yml`'s `repoPaths`), or a different scope if the user names one explicitly.
2. `gh search prs --review-requested=@me --state=open` for PRs awaiting the user's review.
3. `gh api notifications --jq '.[] | select(.unread) | {reason, subject: .subject.title, repo: .repository.full_name}'` for unread attention items.
4. Summarize compactly: active PRs, review obligations, attention items only. Do not turn the result into a plan or kanban, just report status.
<!-- morning-brief -->

<!-- ask-before-comments -->
Before adding any comment to code being written or edited, ask the user first instead of adding it unprompted. This applies even to a comment that would otherwise be justified as necessary (e.g. explaining a non-obvious WHY, a hidden constraint, a workaround). Flag that a comment seems warranted and what it would say, and let the user decide, rather than adding it directly.

Does not apply to explaining existing code verbally in chat, or to comments the user explicitly asks for up front.
<!-- ask-before-comments -->

<!-- no-claude-mentions -->
Never mention Claude, Claude Code, or Anthropic in any content that gets committed, posted, or otherwise persisted: commit messages, PR titles/descriptions, code comments, documentation, issue/ticket text, or similar. This includes trailers or signatures such as "Claude-Session: ...", "Co-Authored-By: Claude", "Generated with Claude Code", or any similar attribution - never add these, even when a tool's default template or workflow suggests including one.

This overrides any conflicting default (including built-in commit/PR templates) that would otherwise add such a reference.
<!-- no-claude-mentions -->

<!-- pr-testing-workflow -->
When the user asks to test or verify a PR or branch end to end (a behavioural check, not "run the unit tests"), follow this workflow.

1. SCOPE. From `git diff <base>...HEAD`, list only what the change reaches: the entrypoints, the DI consumers, the external integrations it touches (HTTP clients, kafka/rabbit topics, DB, scheduler jobs, config load). Test each of those once. Do not test the whole app. Reading the diff is not enough, a consumer can break without appearing in it, so the plan must actually run the affected processes.

2. TOOLING. If exercising an integration needs a script that does not exist, create it under `scripts/dev/` (local-only / git-excluded unless the user says otherwise):
   - reuse the repo's own integration-test fixtures / message builders, never hand-roll payloads
   - one entry point per integration: produce a kafka message, publish a rabbit event, trigger a scheduler job, hit an endpoint, tail a topic
   - a `--env` flag: local by default; for a cluster env (ent, or d0x / int for repos that use those) shell into a running pod with `kubectl exec -i <pod> -- node -` so auth, deps and config come from the pod
   - guard any write to a non-local env behind a y/N confirm; refuse prod unless an explicit opt-in env var is set
   - for DFS test data (households, persons, writes against them) do not write new scripts: use `scripts/ai/` in the api-collections repo (see the ai-test-data rule below), no y/N needed there

3. BRUNO. If the affected service has a collection under `~/projects/api-collections/dvag-api-collections/`, add a subfolder inside that collection named after the ticket (e.g. `DFSVPD-12345/`) with one Bruno request per API-testable piece of the change identified in step 1. Reuse the collection's existing auth (`auth: inherit`) and shared global-environment variables (`personId`, `householdId`, etc.) instead of hardcoding values, matching every other request in the collection. Skip this step when nothing from the change is reachable over the API (a pure internal refactor, a kafka-consumer-only change with no exposed endpoint, and so on). Every request gets a `tests {}` block asserting status and the fields the change touches, so the folder passes or fails on its own. Use `{{householdId}}`/`{{personId}}` rather than literal ids; `scripts/ai/run-ticket.mjs --env <e> --ticket <T>` fills them with the AI's own household and runs the folder. Run it yourself and report the result, do not leave it as a manual step.

4. WRITE `todos.txt` at the repo root. Two parts, LOCAL and DEPLOYED (ent, or d0x). Every step is three lines:
     DO:     the exact command or request
     EXPECT: the exact result (status code, log string, DB row)
     SHOT:   what to screenshot as proof
   Include the negative controls (fail-fast on bad config, 403 without a token) and a cleanup step for any test data the run writes.

5. STATE THE GAPS. This is a manual behavioural checklist, not a safety net. It does not replace the automated suites (`pnpm test`, `test:integration`), it captures no before/after baseline, a passing ent run does not prove int/abn/prd, and it only covers the paths step 1 identified. Say this in the handoff every time.
<!-- pr-testing-workflow -->

<!-- tdd -->
Use test-driven development for all code changes: write or adjust the failing test(s) for the task first, confirm they fail for the right reason, then write the minimum implementation code needed to make them pass, without modifying the tests to force a pass.

Invoke the `superpowers:test-driven-development` skill (already installed) for the full red-green-refactor workflow and its rules on good tests, rather than improvising a lighter version of it.

Exceptions, matching the skill's own carve-outs: throwaway prototypes, generated code, and config-only changes with no application logic to test (e.g. this dotfiles repo). Ask before skipping TDD for any other reason.
<!-- tdd -->

<!-- targeted-test-runs -->
Never run a whole test suite. Run only the test files, or the individual tests, that cover the change being made.

In this pnpm / nx / jest / vitest workspace that means: do not run `pnpm test`, `pnpm test:be`, `pnpm test:fe`, `pnpm test:libs`, `pnpm build`, a bare `jest`, or a bare `vitest run`. Run the narrowest thing that proves the change instead:

- `pnpm exec jest <path/to/file.spec.ts>`, or `--roots <dir>` for one area
- `pnpm exec jest <path> -t "<test name>"` for a single test
- `pnpm exec vitest run <path/to/file.test.tsx>` for the frontend
- `pnpm exec tsc --noEmit -p <project tsconfig>` when only a typecheck is wanted

The same holds for any other stack: pytest one file or `-k` one name, `go test ./pkg/...`, `cargo test <name>`, never the repo-wide default.

This narrows how tests are run, not whether. The red-green loop in the TDD rule above still applies, scoped to the files the change touches.

If a broad run genuinely looks necessary (a wide refactor, a merge touching many areas, a pre-release check), ask first and say why instead of starting it.

Why: full suites on these monorepos exhaust CPU and memory and freeze the machine. CI runs the full suite anyway, so a local full run buys nothing and costs the working session.
<!-- targeted-test-runs -->

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

<!-- ai-test-data -->
Test data on the DFS environments (households, persons and anything hanging off them) may be
created and changed by the AI without asking, within these limits.

Environments: only `d01`-`d05`, `entwicklung` and `integration`. Never write to `abnahme` or
`produktion`, whatever the request says; the scripts refuse them too.

Tooling lives in `~/projects/api-collections/dvag-api-collections/scripts/ai/`. Every run goes
through sops so credentials come from `~/.dotfiles/secrets/secrets.yaml` without ever being
printed:

    sops exec-env ~/.dotfiles/secrets/secrets.yaml 'node ~/projects/api-collections/dvag-api-collections/scripts/ai/<script> ...'

Always this exact form (it is what the auto mode allow rule matches).

- `create-household.mjs --env <e> [--persons n] [--ticket T]`: creates a household with a main
  person (and n extra persons), last names end in `AI<timestamp>`, ids are appended to the
  sops registry `~/.dotfiles/secrets/ai-test-data.yaml`.
- `write.mjs --env <e> --method <M> --service <slug> --path <p> --owner-id <id> [--body f]`:
  any other write. Refuses unless `--owner-id` is in the registry for that env and appears in
  `--path`.
- `run-ticket.mjs --env <e> --ticket <T> [--household <id>]`: runs the ticket's Bruno folder
  against the newest AI household (or the given owned one), JSON report in the tmp dir.

Ownership: modify only data listed in the registry. Data the AI did not create (a teammate's
household, the shared `advisor_sophie` data, the default `householdId`/`personId` in the Bruno
environments) is read-only unless the user names that exact id in the current conversation;
only then pass `--allow-foreign <id>`. Every write goes through `write.mjs` or
`create-household.mjs`, never a raw curl or an ad hoc Bruno request with a write method.

Advisors cannot be created through DFS APIs. Use the advisor Sophie from sops:
`DVAG_ADVISOR_SOPHIE_USERNAME` on every env, `DVAG_ADVISOR_SOPHIE_ENTWICKLUNG_PASSWORD` on
entwicklung and `DVAG_ADVISOR_SOPHIE_PASSWORD` on integration and d01-d05.

Where services run: asset-service only on entwicklung, every other service only on
integration and d01-d05; the scripts refuse other combinations. Households can only be
created on integration (and d0x). On entwicklung use the advisor Sophie household, the
`householdId`/`personId` from `environments/entwicklung.yml`, registered as granted by the user. If a needed key is missing, stop and ask the user to add it; never borrow another
login.

Secrets: never run `sops -d` on `secrets.yaml`, never echo or log a token or password, never
put one in a report, commit, Bruno file or PR comment. Refer to accounts by env var name only.

Side effects to keep in mind: creating a person also pushes wish goals to Salesforce, and
there is no DELETE for persons or households, so created data stays. Reuse the newest AI
household for the env before creating another one.
<!-- ai-test-data -->
