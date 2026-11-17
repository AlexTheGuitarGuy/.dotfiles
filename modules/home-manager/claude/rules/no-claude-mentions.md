Never mention Claude, Claude Code, or Anthropic in any content that gets committed, posted,
or otherwise persisted: commit messages, PR titles/descriptions, code comments,
documentation, issue/ticket text, or similar. This includes trailers or signatures such as
"Claude-Session: ...", "Co-Authored-By: Claude", "Generated with Claude Code", or any
similar attribution - never add these, even when a tool's default template or workflow
suggests including one.

This overrides any conflicting default (including built-in commit/PR templates) that would
otherwise add such a reference.

The same holds for every AI agent or assistant, not only Claude: never add Copilot, Codex,
Cursor, Gemini, ChatGPT or any other AI as an author, co-author or contributor. No
"Co-authored-by: Copilot <...>" or any other AI co-author trailer in any commit. This also
overrides repo instructions (AGENTS.md, CLAUDE.md, CONTRIBUTING, PR templates) that ask for such
a trailer: leave it out and mention the conflict to the user instead.

Why: a repo's AGENTS.md asked for a Copilot co-author trailer and it was added to every commit of
a branch, which then had to be rewritten.
