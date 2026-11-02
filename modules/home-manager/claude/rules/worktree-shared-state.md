A git worktree shares more than the repository with the main checkout and with other worktrees:
the local database, fixed ports, Docker volumes and containers, and caches. Before running
migrations, seeds, resets, a server on a fixed port or `docker compose down -v` in a worktree,
check what else uses that database, port or volume (see prisma-local-db.md and local-servers.md).

Why: a branch without another branch's migration sees drift and offers to reset the shared
database, and a stray server or `down -v` in one worktree breaks or wipes what the others rely on.
