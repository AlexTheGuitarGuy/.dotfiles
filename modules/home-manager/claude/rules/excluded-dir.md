Local-only files (notes, todos, ticket exports, evidence, screenshots, dev scripts, drafts,
walkthroughs) go into `.excluded/` at the root of the current worktree. The global gitignore
ignores `.excluded/` in every repo, so nothing in it is ever committed.

Never edit `.git/info/exclude` (or `.bare/info/exclude`) and never add local-only patterns to a
repo's tracked `.gitignore`. If a local-only file must live at a specific path outside
`.excluded/` (a tool reads it from there), ask the user before adding its pattern to the global
gitignore (`~/.dotfiles/git/.gitignore_global`).

Every worktree has its own `.excluded/`. Removing a worktree deletes it, so move anything worth
keeping first.

Why: per-repo exclude edits piled up unnoticed in every repo; one known folder keeps them tidy
and visible.
