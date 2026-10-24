---
name: skill-distillation
description: Mine past Claude Code sessions for corrections the user keeps making and propose a rule or skill that would have prevented them. Use only when the user asks to improve their rules or skills, or to learn from past sessions. Never edits or installs anything without approval.
---

# Skill Distillation

Turn repeated corrections into one reviewable proposal.

## Sources

Session transcripts live in `~/.claude/projects/<project>/<session>.jsonl`. Read only the user's own
messages first; they carry the corrections:

```bash
jq -r 'select(.type == "user" and (.message.content | type) == "string") | .message.content' <session>.jsonl
```

Scope the search to the projects and time range the user names; otherwise the newest sessions
across all projects. Open the surrounding assistant turns only for the sessions that matter.
Transcripts may contain secrets: never quote tokens, passwords or keys, and keep quotes short.

## Steps

1. Collect corrections: "no", "don't", "I told you", reverts the user asked for, repeated
   reminders. Quote or precisely summarize each with its session.
2. Group them. Keep only behavior that recurs across sessions; drop one-off project constraints.
3. For each group, name the failed default, the preferred behavior, where it applies, the
   exceptions, and the risk of overfitting.
4. Check what already exists: `~/.claude/rules/` (including linked subfolders), `~/.claude/skills/`,
   plugin skills and the memory files. Prefer sharpening an existing rule or skill over adding a
   new one, and say which.
5. Pick the form:
   - **Rule**: a standing constraint that applies everywhere (most corrections). A few lines plus
     a `Why:` line, matching the existing rule files.
   - **Skill**: a multi-step workflow with a clear trigger.
   - **Memory**: a fact about the user or a project, not an instruction.
6. Draft it with precise trigger and anti-trigger wording, the evidence behind it, and where it
   would live. Rules that only matter for one employer or client belong in the linked private
   rules folder, not the public dotfiles.

## Output

One proposal per finding, strongest first: evidence, chosen form, target path, and the draft text.
End with what was considered and dropped as one-off.

Do not write, edit or install any rule, skill or memory until the user approves that proposal.
