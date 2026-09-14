When the user asks to test or verify a PR or branch end to end (a behavioural check,
not "run the unit tests"), follow this workflow.

1. SCOPE. From `git diff <base>...HEAD`, list only what the change reaches: the
   entrypoints, the DI consumers, the external integrations it touches (HTTP clients,
   kafka/rabbit topics, DB, scheduler jobs, config load). Test each of those once. Do
   not test the whole app. Reading the diff is not enough, a consumer can break
   without appearing in it, so the plan must actually run the affected processes.

2. TOOLING. If exercising an integration needs a script that does not exist, create
   it under `scripts/dev/` (local-only / git-excluded unless the user says otherwise):
   - reuse the repo's own integration-test fixtures / message builders, never
     hand-roll payloads
   - one entry point per integration: produce a kafka message, publish a rabbit
     event, trigger a scheduler job, hit an endpoint, tail a topic
   - a `--env` flag: local by default; for a cluster env (ent, or d0x / int for
     repos that use those) shell into a running pod with
     `kubectl exec -i <pod> -- node -` so auth, deps and config come from the pod
   - guard any write to a non-local env behind a y/N confirm; refuse prod unless an
     explicit opt-in env var is set

3. WRITE `todos.txt` at the repo root. Two parts, LOCAL and DEPLOYED (ent, or d0x).
   Every step is three lines:
     DO:     the exact command or request
     EXPECT: the exact result (status code, log string, DB row)
     SHOT:   what to screenshot as proof
   Include the negative controls (fail-fast on bad config, 403 without a token) and a
   cleanup step for any test data the run writes.

4. STATE THE GAPS. This is a manual behavioural checklist, not a safety net. It does
   not replace the automated suites (`pnpm test`, `test:integration`), it captures no
   before/after baseline, a passing ent run does not prove int/abn/prd, and it only
   covers the paths step 1 identified. Say this in the handoff every time.
