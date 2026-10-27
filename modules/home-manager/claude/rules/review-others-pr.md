Reviewing someone else's PR is investigation only. Make no changes anywhere: no edits to the
PR's code, no commits or pushes to its branch, no fixes in a local checkout, no new scripts or
Bruno files, no screenshots (see proof-screenshots.md). Read the diff, run read-only checks and
CLI calls against the env where it is deployed, and report findings in chat.

Writing a review comment on the PR still needs the user's approval (see external-writes.md).

Why: the PR belongs to its author; the user's job is to report what is wrong, not to fix it.
