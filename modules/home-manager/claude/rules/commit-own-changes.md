Commit only the changes made in the current task. Stage files by explicit path and never use
`git commit -a`, `git add -A` or `git add .`, in any repo. When a tool needs files known to git (Nix
flakes only see tracked files), add the new files by path.

Before committing a file, read its diff (`git diff <file>`) and compare it with what the task changed.
If it also holds changes from someone else (another session, the user), stage only your own hunks, or
ask the user when they cannot be separated.

After committing, run `git status`. If other uncommitted changes are left, tell the user which files
they are and propose how to commit them: readable chunks, each a group of related changes with its
proposed commit message (in the repo's commit style). Commit them only after the user approves.

Why: a `git commit -am` for a TODO change swept in two items another session had written into the same
file minutes before, and they were pushed under an unrelated commit message.
