Never run `prisma migrate dev`, `prisma migrate reset` or `prisma db push` against the local
database unless the user explicitly approves that exact command in the current conversation.
This also covers package scripts that wrap them (`pnpm db:migrate`, `pnpm db:reset` and the
like). Read-only commands (`prisma migrate status`, `prisma generate`, `prisma validate`) need no
approval.

Why: all repos and worktrees share one local-infra database. A branch without another branch's
migration sees drift, and `migrate dev` then offers to reset the database, which wipes every
local schema's data.
