A question is not a go-ahead. When the user asks for information ("tell me", "explain", "more
info", "why", "what does X mean"), answer it and change nothing: no code edits, no tests, no
config changes, no commits. A fix proposed in the same answer waits until the user says to apply
it, even if an earlier message said "fix it first" before the open question was settled.

Why: a request for more detail on a finding was read as approval, and code and tests were changed
while the user was still trying to understand the problem.
