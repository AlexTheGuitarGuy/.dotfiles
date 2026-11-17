In TypeScript code you write or change, never use `any` and never use type assertions: no
`as T`, `as unknown as T`, `as never`, `<T>value`, and no non-null assertion `value!` to silence
the compiler. This holds for production code, tests, mocks and test helpers alike.

Use the type system instead: precise types and generics, function overloads, type predicates
(`value is T`) after a real check, schema parsing (zod) for JSON and other external input,
`unknown` plus narrowing, and typed mocks (`jest-mock-extended` `mock<T>()`) instead of casting
partial objects. `as const` is fine. Assertions that already exist on the base branch stay
untouched unless the user asks; only the code you add or change must be free of them.

Never switch off `no-explicit-any` or similar lint rules (in the eslint config or with
`eslint-disable` comments) to get code through. Passing lint is not proof: many repos do not
lint for assertions, so check your own diff (`git diff <base>... | grep -E ' as |: any|<any>'`).
Lint the whole affected project (for example `nx run <project>:lint`), not only the changed
files: rules that span files (module boundaries, lazy-loaded libraries) only fail in the files
that did not change.

Before writing code in a repo, read its instruction files and follow them: `AGENTS.md`,
`CLAUDE.md`, `.github/copilot-instructions.md`, `.github/instructions/*.md`, `CONTRIBUTING.md`.
Rules written for another assistant (Copilot) apply to you as well.

Why: a Prisma migration PR added `any` in tests, switched off the lint rule, and added 30 type
assertions, although the repo's `.github/instructions/general_rules.instructions.md` forbids
all of it. Nobody had read that file, and lint stayed green.
