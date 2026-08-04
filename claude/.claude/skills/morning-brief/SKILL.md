---
name: morning-brief
description: Use for daily triage, "what needs attention", or a compact current-work GitHub summary — active PRs authored by the user, PRs awaiting their review, and unread notifications. Trigger on explicit ask, never automatically at session start.
---

# Morning Brief

Compact GitHub work status: active PRs, review obligations, attention items. Not
a plan, not a kanban, just status.

## Steps

1. `gh pr list --author=@me --state=open` for PRs the user has open, scoped to
   `dvag/*` orgs by default (matches the repo mapping already in
   `gh-dash/.config/gh-dash/config.yml`'s `repoPaths`), or a different scope if
   the user names one explicitly.
2. `gh search prs --review-requested=@me --state=open` for PRs awaiting the
   user's review.
3. `gh api notifications --jq '.[] | select(.unread) | {reason, subject: .subject.title, repo: .repository.full_name}'`
   for unread attention items.
4. Summarize compactly: active PRs, review obligations, attention items only.
   Do not turn the result into a plan or kanban, just report status.
