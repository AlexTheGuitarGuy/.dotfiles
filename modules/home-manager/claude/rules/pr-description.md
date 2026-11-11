PR descriptions follow two content rules on top of the repo's template.

- PR Type: tick exactly one checkbox, the one matching the conventional-commit type of the PR
  title (`feat` is New Feature). Never tick extra boxes for docs or tests that come with it.
- No local machine details: nothing about NixOS, direnv, local engine or toolchain workarounds,
  the agent's secret or registry setup, or other quirks of this workstation. Commands in the PR
  are the ones a teammate runs on a standard setup (the repo's package script, not a
  workstation-specific wrapper). The same holds for committed docs (README, runbooks).

Keep PR text to the essentials: what changed, why, how to run or verify it, what is not covered.
No restating the diff, no long prose.

Assign the user to every PR created for them: `gh pr create --assignee @me` (or `gh pr edit <n>
--add-assignee @me` afterwards).

To request or re-request a Copilot review, use GraphQL; `gh pr edit --add-reviewer @copilot` and
the REST `requested_reviewers` call silently do nothing:

    gh api graphql -f query='mutation{requestReviewsByLogin(input:{pullRequestId:"<PR node id>",
      botLogins:["copilot-pull-request-reviewer"],union:true}){pullRequest{id}}}'

Why: the user corrected a PR that ticked three PR types, told reviewers to use direnv on NixOS,
and was too verbose. A PR was also opened without the user as assignee, who had to add it by hand.
