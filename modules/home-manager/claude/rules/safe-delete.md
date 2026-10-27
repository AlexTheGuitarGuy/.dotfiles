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
