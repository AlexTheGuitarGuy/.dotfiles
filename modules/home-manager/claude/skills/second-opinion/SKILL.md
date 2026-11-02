---
name: second-opinion
description: Use only when the user asks for a second opinion on an architecture decision or plan. Writes an evidence packet, then asks a Claude Opus subagent and an OpenCode GPT run for independent opinions, each blind to the other, and compares them.
---

# Second opinion

Independent opinions on one architecture or plan question, after the evidence is gathered. Never run
this unprompted and never for a code diff (use a review skill for that).

## Steps

1. Gather the evidence first. Write a packet file in the worktree's `.excluded/` (for example
   `.excluded/second-opinion/<topic>.md`) with: the question, the requirements, the options
   considered, the relevant file paths and excerpts, constraints, and what is already decided. State
   facts and sources. Do not state your own preferred answer, so neither reviewer is steered.

2. Start both reviewers in the same turn, in parallel, and give each only the packet path and the
   prompt below. Neither sees the other's answer or your leaning.

   - Claude: an Agent call with `model: "opus"` and a read-only brief (read the packet and the
     files it names, change nothing).
   - OpenCode, via Bash with a read-only agent:

         opencode run -m github-copilot/gpt-6.1-sol --agent explore "<prompt with the packet path>"

     `explore` is read-only in the OpenCode config. The model id must exist for the
     `github-copilot` provider (`opencode models | grep github-copilot`); if it does not, say so and
     continue with the Claude opinion only.

3. Prompt for both:

   > Review only the architecture or plan described in the packet at <path>, after reading the files
   > it cites. Evaluate alternatives against requirements, existing boundaries, failure modes and
   > operational cost, using cited sources. Return a recommendation, the strongest objection,
   > assumptions that could change it, and the decisive check, in at most 700 words. Distinguish
   > evidence from inference. Do not invent problems to disagree, and do not claim agreement
   > without independent analysis. Do not implement anything, review a diff, or modify files.

4. Compare. Report where the two agree, where they differ, which assumptions each depends on, and
   the one check that would settle a disagreement. Give your own recommendation last, marked as
   yours. Do not average the opinions.
