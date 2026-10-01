---
name: review-live
description: Use when reviewing a DFS PR, branch or own change and the code is deployed somewhere. Calls the affected endpoints and flows on the env where the change runs (preferred over local) and reports every issue with exact reproduction steps.
---

# Review Live

Static review finds what the diff says. This finds what the deployed code does. Run it
alongside the normal code review, not instead of it. Follow ai-test-data.md for every call.

## 1. Scope

From the diff, list what a caller can reach: REST/GraphQL endpoints, consumers of changed
DTOs, flows that end in the changed code. Use gitnexus (`detect_changes`, `route_map`,
`api_impact`, `impact`) to find callers the diff does not show. One line per item:
method, path, what changed about it.

## 2. Find where it is deployed

Local only if nothing else runs the change, and say so.

    gh pr view <n> --json headRefName,headRefOid
    gh api "repos/dvag/<repo>/deployments?ref=<headRefName>" \
      -q '.[] | {id, environment, sha, created_at}'
    gh api repos/dvag/<repo>/deployments/<id>/statuses -q '.[0].state'

An env counts only if the latest `success` deployment there has `sha == headRefOid` and no
newer successful deployment to the same env came from another ref
(`gh api "repos/dvag/<repo>/deployments?environment=<env>"`). For dfs-asset-service,
`dfs-entwicklung` is the target (`entwicklung`); for every other repo ignore `dfs-entwicklung`
(those services do not run there) and use `integration` or d0x only. If the head sha is not deployed anywhere, say which older sha is and whether
the diff since then matters. d0x deploys do not appear in the deployments API: if nothing
matches, ask the user which d0N runs the branch.

A `failure`/`error` status does not prove the pod is old (the status can come from a later
pipeline step). When the head sha only has failed deployments, report that, and ask the user
before treating the env as running the change.

## 3. Exercise it

- Data: reuse the newest AI household for the env, create one with `create-household.mjs`
  only if none exists or the flow needs a fresh one.
- Reads and writes: the ticket's Bruno folder via `run-ticket.mjs` when one exists, otherwise
  `write.mjs` for writes and `bru run <request.bru> --global-env <env>` for reads.
- Per item from step 1: the happy path, one invalid input, and a call without token or with
  a foreign household where the change touches authorization.

## 4. Report

For every issue, a block the reader can replay without asking anything:

    ISSUE:    one line
    ENV:      entwicklung @ <sha>   (deployment <id>)
    DATA:     household <id> (AI-created) | advisor Sophie via DVAG_ADVISOR_SOPHIE_*
    REPRO:    exact command (run-ticket / write.mjs / bru run), request body inline
    EXPECTED: status + fields
    ACTUAL:   status + response excerpt (no tokens)
    CAUSE:    file:line in the diff, if known

Then list what was exercised and passed, and what could not be reached and why. Posting any
of this to the PR still needs the user's approval (external-writes rule).
