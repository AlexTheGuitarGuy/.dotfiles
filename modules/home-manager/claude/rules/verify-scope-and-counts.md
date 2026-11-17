Before stating which code a branch introduced, how many places are affected, or that something is
"ours", establish it exactly and check it a second way.

- Attribution by blame, not by diff counts. "Added in the diff" or "more in this file than on the
  base" is not "written on this branch": moved or renamed files and rewritten lines carry old code
  along. Use `git blame -w -M -C -C` on the branch head and keep only lines whose commit is in
  `git rev-list <base>..HEAD`. Lines from commits already on the base belong to other people and
  stay untouched unless the user says otherwise.
- Test a search pattern before trusting its count. Run it against a few known positives and
  negatives (cases you have already seen in this code) and widen it until it finds them all. Say
  which pattern was used and what it cannot catch.
- No number or yes/no claim as fact from a single method. Confirm it a second, independent way
  (a different pattern, blame next to the diff, a manual sample). If they disagree, say so and
  resolve it before acting on it.
- Subagents get exact item lists (`file:line` plus the line text), never a rule of thumb like
  "remove as many as the base did not have". Lines that must not change get their own list too.
  Check the result against both lists before committing or pushing.

Why: a review counted 30 type assertions a PR had added with a too narrow pattern (the real
number was about 104), then told a subagent to remove "as many per file as main did not have",
which would also have hit 31 lines from other people's commits that blame showed were already on
main.
