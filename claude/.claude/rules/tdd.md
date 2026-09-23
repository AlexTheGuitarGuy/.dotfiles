Use test-driven development for all code changes: write or adjust the failing test(s) for
the task first, confirm they fail for the right reason, then write the minimum
implementation code needed to make them pass, without modifying the tests to force a pass.

Invoke the `superpowers:test-driven-development` skill (already installed) for the full
red-green-refactor workflow and its rules on good tests, rather than improvising a lighter
version of it.

Exceptions, matching the skill's own carve-outs: throwaway prototypes, generated code, and
config-only changes with no application logic to test (e.g. this dotfiles repo). Ask before
skipping TDD for any other reason.
