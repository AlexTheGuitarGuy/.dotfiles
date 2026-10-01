Never run a whole test suite. Run only the test files, or the individual tests, that cover the
change being made.

In this pnpm / nx / jest / vitest workspace that means: do not run `pnpm test`, `pnpm test:be`,
`pnpm test:fe`, `pnpm test:libs`, `pnpm build`, a bare `jest`, or a bare `vitest run`. Run the
narrowest thing that proves the change instead:

- `pnpm exec jest <path/to/file.spec.ts>`, or `--roots <dir>` for one area
- `pnpm exec jest <path> -t "<test name>"` for a single test
- `pnpm exec vitest run <path/to/file.test.tsx>` for the frontend
- `pnpm exec tsc --noEmit -p <project tsconfig>` when only a typecheck is wanted

The same holds for any other stack: pytest one file or `-k` one name, `go test ./pkg/...`,
`cargo test <name>`, never the repo-wide default.

This narrows how tests are run, not whether. The red-green loop in tdd.md still applies, scoped
to the files the change touches.

If a broad run genuinely looks necessary (a wide refactor, a merge touching many areas, a
pre-release check), ask first and say why instead of starting it.

Why: full suites on these monorepos exhaust CPU and memory and freeze the machine. CI runs the
full suite anyway, so a local full run buys nothing and costs the working session.
