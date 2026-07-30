Before doing open-ended research (e.g. "what is X"), across any project:

1. Use any CLI needed (az, kubectl, helm, git, etc.) without asking first, for read-only exploration.
2. Look anywhere in the relevant workspace/repos. If a repo isn't on the latest commit, local changes may be stashed and main checked out and pulled, but ask permission first before that checkout/pull/stash step.
3. Use any MCP needed for the research.
4. Never touch prod or staging/UAT environments without asking first. Any non-readonly operation (write/delete/exec/restart) on a pod, database, or similar resource, in any environment, also needs to be asked about first.

This is standing permission for read-only exploration (grep, CLI calls, MCP queries) so research can be self-served rather than asking each time, with carve-outs for repo state mutation and for anything touching real infrastructure.
