Every commit on a branch must use the branch/PR's conventional-commit type, not the literal
nature of that individual change. If the branch/PR is `feat(...)`, every commit on it is
`feat(...)` too, even ones that would look like a `fix` or `docs` commit in isolation (bug fix,
doc update, test-only change). Never mix `fix(...)`, `docs(...)`, etc. into a `feat(...)` branch,
or vice versa.

Before writing any commit message, check the branch name, PR title, or the branch's first commit
for its conventional-commit type, and prefix every subsequent commit on that branch with that
same type.
