Commit only the changes made in the current task. Stage files by explicit path and never use
`git commit -a`, `git add -A` or `git add .` in a repo that other sessions also write to (the dotfiles
and work repos, shared TODO files, memory). Exception: a repo whose own instructions say to stage
everything, and only after the check below.

Before committing a file, read its diff (`git diff <file>`) and compare it with what the task changed.
If it holds changes from someone else (another session, the user), leave them out: stage only your own
hunks, or ask the user when they cannot be separated. Say in the report what was left uncommitted.

Why: a `git commit -am` for a TODO change swept in two items another session had written into the same
file minutes before, and they were pushed under an unrelated commit message.
