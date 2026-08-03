<!-- context7 -->
Use Context7 MCP to fetch current documentation whenever the user asks about a library, framework, SDK, API, CLI tool, or cloud service -- even well-known ones like React, Next.js, Prisma, Express, Tailwind, Django, or Spring Boot. This includes API syntax, configuration, version migration, library-specific debugging, setup instructions, and CLI tool usage. Use even when you think you know the answer -- your training data may not reflect recent changes. Prefer this over web search for library docs.

Do not use for: refactoring, writing scripts from scratch, debugging business logic, code review, or general programming concepts.

## Steps

1. Always start with `resolve-library-id` using the library name and the user's question, unless the user provides an exact library ID in `/org/project` format
2. Pick the best match (ID format: `/org/project`) by: exact name match, description relevance, code snippet count, source reputation (High/Medium preferred), and benchmark score (higher is better). If results don't look right, try alternate names or queries (e.g., "next.js" not "nextjs", or rephrase the question). Use version-specific IDs when the user mentions a version
3. `query-docs` with the selected library ID and the user's full question (not single words)
4. Answer using the fetched docs
<!-- context7 -->

<!-- no-em-dash -->
Never use the em dash (—) or en dash (–) in written text: chat responses, code comments, commit messages, documentation, or any other prose, in any project.

Use a comma, period, semicolon, parentheses, or a connecting word (and, but, since, because) instead.
<!-- no-em-dash -->

<!-- no-em-dash-scope -->
The no-em-dash rule (see above) governs only text I author myself.

Never retroactively "clean up" an em dash or en dash in existing file content that someone else wrote, even while substantially rewriting other parts of the same file for an unrelated task. When editing a file for a specific task, only avoid the dashes in the new text being added; do not scan the rest of the file for dashes to fix as a side effect. If genuine em-dash cleanup of existing prose is wanted, that must be an explicit, separate request from the user.
<!-- no-em-dash-scope -->

<!-- external-writes -->
Any write or modification to a system visible to other people (Jira tickets, Confluence pages, GitHub/PR comments, Slack messages, emails, or similar) requires explicit approval from the user before posting, every time, no standing permission.

Before posting, the content must:
1. Read as authored by the user, not by an AI. No AI-tell phrasing, no signature or disclosure that it was AI-written.
2. Follow the no-em-dash rule above, since this is user-facing/external-facing text.

Read-only lookups (viewing tickets, pages, threads) need no approval. Any create/update/comment/transition action always needs a review-and-approve step first, drafted in plain prose without dashes or AI-sounding phrasing.

Why: this kind of content is visible to teammates and represents the user's own voice/authorship, unlike local file edits or memory notes.
<!-- external-writes -->

<!-- research-permissions -->
Before doing open-ended research (e.g. "what is X"), across any project:

1. Use any CLI needed (az, kubectl, helm, git, etc.) without asking first, for read-only exploration.
2. Look anywhere in the relevant workspace/repos. If a repo isn't on the latest commit, local changes may be stashed and main checked out and pulled, but ask permission first before that checkout/pull/stash step.
3. Use any MCP needed for the research.
4. Never touch prod or staging/UAT environments without asking first. Any non-readonly operation (write/delete/exec/restart) on a pod, database, or similar resource, in any environment, also needs to be asked about first.

This is standing permission for read-only exploration (grep, CLI calls, MCP queries) so research can be self-served rather than asking each time, with carve-outs for repo state mutation and for anything touching real infrastructure.
<!-- research-permissions -->
