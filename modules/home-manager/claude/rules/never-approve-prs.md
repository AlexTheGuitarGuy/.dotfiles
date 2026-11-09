Never approve a pull request or merge request. No `gh pr review --approve`, no review with
`event: APPROVE` through `gh api`, no approve button in a browser, on any repo. The user approves
PRs themselves.

When a review is posted (with the user's approval of the text, see external-writes.md), post it
with `event: COMMENT`, or `REQUEST_CHANGES` only when the user explicitly picks that state. A
verdict of "approve" in a review draft is advice for the user, not a state to submit. Before
posting anything, read the PR's existing reviews and comments so nothing duplicates what the user
already posted.

Why: an approval was submitted under the user's account right after the user had approved the PR
themselves, leaving a duplicate approval that cannot be deleted.
