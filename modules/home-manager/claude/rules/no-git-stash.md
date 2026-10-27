Never use `git stash`. To set local changes aside, copy the files to `.excluded/backup/` in the
worktree (see excluded-dir.md), verify the copy, then restore or discard them in the tree.

Why: a stash and pop during a session broke the working tree and had to be recovered by hand.
